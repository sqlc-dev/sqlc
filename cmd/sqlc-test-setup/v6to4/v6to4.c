// v6to4 is an LD_PRELOAD library that lets a program which insists on IPv6
// sockets run on a kernel built without IPv6, where socket(AF_INET6, ...)
// fails with EAFNOSUPPORT. sqlc-test-setup preloads it into Spanner Omni on
// such kernels only, since Omni's services open IPv6 sockets and abort when
// they cannot; see omni_ipv6.go.
//
// An AF_INET6 socket is created as an AF_INET one and remembered. Addresses
// passed in are rewritten from sockaddr_in6 to sockaddr_in (:: to 0.0.0.0,
// ::1 to 127.0.0.1, ::ffff:a.b.c.d to a.b.c.d), and addresses handed back
// are rewritten the other way, so the program goes on seeing IPv6. An
// address with no IPv4 equivalent fails with EADDRNOTAVAIL.
//
// Omni runs on its own glibc 2.27, so the library calls the kernel through
// syscall() rather than looking up the real functions, and must not need a
// glibc symbol newer than that: build it without _FORTIFY_SOURCE or the
// stack protector.
#define _GNU_SOURCE
#include <errno.h>
#include <netinet/in.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/syscall.h>
#include <unistd.h>

#define MAXFD 65536
static unsigned char fake[MAXFD];

static int isfake(int fd) { return fd >= 0 && fd < MAXFD && fake[fd]; }
static void mark(int fd, int v) { if (fd >= 0 && fd < MAXFD) fake[fd] = v; }

// to4 rewrites a sockaddr_in6 as a sockaddr_in in out; it returns the length to
// pass on, or -1 when the address has no IPv4 equivalent.
static socklen_t to4(const struct sockaddr *a, socklen_t len, struct sockaddr_in *out) {
	const struct sockaddr_in6 *s6 = (const struct sockaddr_in6 *)a;
	static const unsigned char any[16] = {0};
	static const unsigned char lo[16] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1};
	memset(out, 0, sizeof *out);
	out->sin_family = AF_INET;
	out->sin_port = s6->sin6_port;
	const unsigned char *b = s6->sin6_addr.s6_addr;
	if (memcmp(b, any, 16) == 0) out->sin_addr.s_addr = htonl(INADDR_ANY);
	else if (memcmp(b, lo, 16) == 0) out->sin_addr.s_addr = htonl(INADDR_LOOPBACK);
	else if (memcmp(b, any, 10) == 0 && b[10] == 0xff && b[11] == 0xff) memcpy(&out->sin_addr, b + 12, 4);
	else return (socklen_t)-1;
	(void)len;
	return sizeof *out;
}

// to6 rewrites the sockaddr_in the kernel wrote into a v4-mapped sockaddr_in6.
static void to6(struct sockaddr *a, socklen_t *len, socklen_t cap) {
	if (!a || !len || *len < sizeof(struct sockaddr_in) || a->sa_family != AF_INET) return;
	struct sockaddr_in s4;
	memcpy(&s4, a, sizeof s4);
	struct sockaddr_in6 s6;
	memset(&s6, 0, sizeof s6);
	s6.sin6_family = AF_INET6;
	s6.sin6_port = s4.sin_port;
	if (s4.sin_addr.s_addr == htonl(INADDR_LOOPBACK)) s6.sin6_addr.s6_addr[15] = 1;
	else if (s4.sin_addr.s_addr != htonl(INADDR_ANY)) {
		s6.sin6_addr.s6_addr[10] = 0xff;
		s6.sin6_addr.s6_addr[11] = 0xff;
		memcpy(s6.sin6_addr.s6_addr + 12, &s4.sin_addr, 4);
	}
	memcpy(a, &s6, cap < sizeof s6 ? cap : sizeof s6);
	*len = sizeof s6;
}

int socket(int domain, int type, int protocol) {
	int fake6 = domain == AF_INET6;
	int fd = syscall(SYS_socket, fake6 ? AF_INET : domain, type, protocol);
	if (fd >= 0) mark(fd, fake6);
	return fd;
}

int close(int fd) {
	mark(fd, 0);
	return syscall(SYS_close, fd);
}

#define WITH4(fd, addr, len, call4, call)                                      \
	if ((addr) && (addr)->sa_family == AF_INET6 && isfake(fd)) {              \
		struct sockaddr_in s4;                                                 \
		socklen_t l4 = to4((addr), (len), &s4);                                \
		if (l4 == (socklen_t)-1) { errno = EADDRNOTAVAIL; return -1; }        \
		return call4;                                                          \
	}                                                                          \
	return call;

int bind(int fd, const struct sockaddr *addr, socklen_t len) {
	WITH4(fd, addr, len, syscall(SYS_bind, fd, &s4, l4), syscall(SYS_bind, fd, addr, len))
}

int connect(int fd, const struct sockaddr *addr, socklen_t len) {
	WITH4(fd, addr, len, syscall(SYS_connect, fd, &s4, l4), syscall(SYS_connect, fd, addr, len))
}

ssize_t sendto(int fd, const void *buf, size_t n, int flags, const struct sockaddr *addr, socklen_t len) {
	WITH4(fd, addr, len, syscall(SYS_sendto, fd, buf, n, flags, &s4, l4),
	      syscall(SYS_sendto, fd, buf, n, flags, addr, len))
}

ssize_t sendmsg(int fd, const struct msghdr *msg, int flags) {
	if (msg && msg->msg_name && isfake(fd) && ((struct sockaddr *)msg->msg_name)->sa_family == AF_INET6) {
		struct sockaddr_in s4;
		socklen_t l4 = to4(msg->msg_name, msg->msg_namelen, &s4);
		if (l4 == (socklen_t)-1) { errno = EADDRNOTAVAIL; return -1; }
		struct msghdr m = *msg;
		m.msg_name = &s4;
		m.msg_namelen = l4;
		return syscall(SYS_sendmsg, fd, &m, flags);
	}
	return syscall(SYS_sendmsg, fd, msg, flags);
}

#define BACK6(fd, ret, addr, lenp)                                             \
	do {                                                                       \
		if ((ret) >= 0 && isfake(fd) && (addr) && (lenp)) to6((addr), (lenp), cap); \
	} while (0)

int getsockname(int fd, struct sockaddr *addr, socklen_t *len) {
	socklen_t cap = len ? *len : 0;
	struct sockaddr_in6 tmp; socklen_t tl = sizeof tmp;
	if (!isfake(fd)) return syscall(SYS_getsockname, fd, addr, len);
	int r = syscall(SYS_getsockname, fd, &tmp, &tl);
	if (r < 0) return r;
	to6((struct sockaddr *)&tmp, &tl, sizeof tmp);
	memcpy(addr, &tmp, cap < tl ? cap : tl);
	*len = tl;
	return r;
}

int getpeername(int fd, struct sockaddr *addr, socklen_t *len) {
	socklen_t cap = len ? *len : 0;
	struct sockaddr_in6 tmp; socklen_t tl = sizeof tmp;
	if (!isfake(fd)) return syscall(SYS_getpeername, fd, addr, len);
	int r = syscall(SYS_getpeername, fd, &tmp, &tl);
	if (r < 0) return r;
	to6((struct sockaddr *)&tmp, &tl, sizeof tmp);
	memcpy(addr, &tmp, cap < tl ? cap : tl);
	*len = tl;
	return r;
}

static int acc(int fd, struct sockaddr *addr, socklen_t *len, int flags) {
	if (!isfake(fd)) return syscall(SYS_accept4, fd, addr, len, flags);
	struct sockaddr_in6 tmp; socklen_t tl = sizeof tmp;
	int r = syscall(SYS_accept4, fd, &tmp, &tl, flags);
	if (r < 0) return r;
	mark(r, 1);
	if (addr && len) {
		socklen_t cap = *len;
		to6((struct sockaddr *)&tmp, &tl, sizeof tmp);
		memcpy(addr, &tmp, cap < tl ? cap : tl);
		*len = tl;
	}
	return r;
}
int accept(int fd, struct sockaddr *addr, socklen_t *len) { return acc(fd, addr, len, 0); }
int accept4(int fd, struct sockaddr *addr, socklen_t *len, int flags) { return acc(fd, addr, len, flags); }

ssize_t recvfrom(int fd, void *buf, size_t n, int flags, struct sockaddr *addr, socklen_t *len) {
	socklen_t cap = len ? *len : 0;
	ssize_t r = syscall(SYS_recvfrom, fd, buf, n, flags, addr, len);
	BACK6(fd, r, addr, len);
	return r;
}

ssize_t recvmsg(int fd, struct msghdr *msg, int flags) {
	socklen_t cap = msg ? msg->msg_namelen : 0;
	ssize_t r = syscall(SYS_recvmsg, fd, msg, flags);
	if (r >= 0 && msg && msg->msg_name) BACK6(fd, r, (struct sockaddr *)msg->msg_name, &msg->msg_namelen);
	return r;
}

int setsockopt(int fd, int level, int name, const void *val, socklen_t len) {
	if (isfake(fd) && level == IPPROTO_IPV6) {
		if (name == IPV6_TCLASS) return syscall(SYS_setsockopt, fd, IPPROTO_IP, IP_TOS, val, len);
		return 0; // IPV6_V6ONLY and the rest have no IPv4 meaning
	}
	return syscall(SYS_setsockopt, fd, level, name, val, len);
}

int getsockopt(int fd, int level, int name, void *val, socklen_t *len) {
	if (isfake(fd) && level == SOL_SOCKET && name == SO_DOMAIN && val && len && *len >= sizeof(int)) {
		*(int *)val = AF_INET6;
		*len = sizeof(int);
		return 0;
	}
	if (isfake(fd) && level == IPPROTO_IPV6) {
		if (val && len && *len >= sizeof(int)) { *(int *)val = 0; *len = sizeof(int); }
		return 0;
	}
	return syscall(SYS_getsockopt, fd, level, name, val, len);
}
