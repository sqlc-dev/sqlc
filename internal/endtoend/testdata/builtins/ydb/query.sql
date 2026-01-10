-- Basic functions with concrete return types
-- name: FLength :one
SELECT LENGTH("hello");
-- name: FLen :one
SELECT LEN("world");
-- name: FSubstring2 :one
SELECT Substring("abcdef", 2);
-- name: FSubstring3 :one
SELECT Substring("abcdef", 2, 3);
-- name: FFind2 :one
SELECT Find("abcdef", "cd");
-- name: FFind3 :one
SELECT Find("abcdef", "c", 3);
-- name: FRFind2 :one
SELECT RFind("ababa", "ba");
-- name: FStartsWith :one
SELECT StartsWith("abcdef", "abc");
-- name: FEndsWith :one
SELECT EndsWith("abcdef", "def");
-- name: FIf2 :one
SELECT IF(true, 1);
-- name: FIf3 :one
SELECT IF(false, 1, 2);
-- name: FCurrentUtcDate :one
SELECT CurrentUtcDate();
-- name: FCurrentUtcDatetime :one
SELECT CurrentUtcDatetime();
-- name: FCurrentUtcTimestamp :one
SELECT CurrentUtcTimestamp();
-- name: FVersion :one
SELECT Version();
-- name: FToBytes :one
SELECT ToBytes("abc");
-- name: FByteAt :one
SELECT ByteAt("abc", 1);
-- name: FTestBit :one
SELECT TestBit("a", 0);
-- name: FClearBit :one
SELECT ClearBit("a", 0);
-- name: FSetBit :one
SELECT SetBit("a", 0);
-- name: FFlipBit :one
SELECT FlipBit("a", 0);

-- Math functions with concrete return types
-- name: FMathPi :one
SELECT Math::Pi();
-- name: FMathE :one
SELECT Math::E();
-- name: FMathAbs :one
SELECT Math::Abs(-5.5);
-- name: FMathAcos :one
SELECT Math::Acos(0.5);
-- name: FMathAsin :one
SELECT Math::Asin(0.5);
-- name: FMathAtan :one
SELECT Math::Atan(1.0);
-- name: FMathCbrt :one
SELECT Math::Cbrt(27.0);
-- name: FMathCeil :one
SELECT Math::Ceil(4.2);
-- name: FMathCos :one
SELECT Math::Cos(0.0);
-- name: FMathExp :one
SELECT Math::Exp(1.0);
-- name: FMathFloor :one
SELECT Math::Floor(4.8);
-- name: FMathLog :one
SELECT Math::Log(2.718281828);
-- name: FMathLog2 :one
SELECT Math::Log2(8.0);
-- name: FMathLog10 :one
SELECT Math::Log10(100.0);
-- name: FMathRound :one
SELECT Math::Round(4.6);
-- name: FMathRound2 :one
SELECT Math::Round(4.567, 2);
-- name: FMathSin :one
SELECT Math::Sin(0.0);
-- name: FMathSqrt :one
SELECT Math::Sqrt(16.0);
-- name: FMathTan :one
SELECT Math::Tan(0.0);
-- name: FMathTrunc :one
SELECT Math::Trunc(4.9);
-- name: FMathAtan2 :one
SELECT Math::Atan2(1.0, 1.0);
-- name: FMathPow :one
SELECT Math::Pow(2.0, 3.0);
-- name: FMathHypot :one
SELECT Math::Hypot(3.0, 4.0);
-- name: FMathFmod :one
SELECT Math::Fmod(10.5, 3.0);
-- name: FMathIsinf :one
SELECT Math::Isinf(1.0/0.0);
-- name: FMathIsnan :one
SELECT Math::Isnan(0.0/0.0);
-- name: FMathIsfinite :one
SELECT Math::Isfinite(5.0);
-- name: FMathFuzzyequals :one
SELECT Math::Fuzzyequals(1.0, 1.0001);
-- name: FMathMod :one
SELECT Math::Mod(10, 3);
-- name: FMathRem :one
SELECT Math::Rem(10, 3);

-- DateTime functions with concrete return types
-- name: FDateTimeGetyear :one
SELECT DateTime::Getyear(CurrentUtcDate());
-- name: FDateTimeGetmonth :one
SELECT DateTime::Getmonth(CurrentUtcDate());
-- name: FDateTimeGetdayofmonth :one
SELECT DateTime::Getdayofmonth(CurrentUtcDate());
-- name: FDateTimeGethour :one
SELECT DateTime::Gethour(CurrentUtcDatetime());
-- name: FDateTimeGetminute :one
SELECT DateTime::Getminute(CurrentUtcDatetime());
-- name: FDateTimeGetsecond :one
SELECT DateTime::Getsecond(CurrentUtcDatetime());
-- name: FDateTimeFromseconds :one
SELECT DateTime::Fromseconds(1640995200);
-- name: FDateTimeFrommilliseconds :one
SELECT DateTime::Frommilliseconds(1640995200000);
-- name: FDateTimeIntervalfromdays :one
SELECT DateTime::Intervalfromdays(7);

-- Unicode functions with concrete return types
-- name: FUnicodeIsutf :one
SELECT Unicode::Isutf("hello");
-- name: FUnicodeGetlength :one
SELECT Unicode::Getlength("你好");
-- name: FUnicodeFind :one
SELECT Unicode::Find("hello", "ll");
-- name: FUnicodeRfind :one
SELECT Unicode::Rfind("hello", "l");
-- name: FUnicodeSubstring :one
SELECT Unicode::Substring("hello", 1, 3);
-- name: FUnicodeNormalize :one
SELECT Unicode::Normalize("café");
-- name: FUnicodeTolower :one
SELECT Unicode::Tolower("HELLO");
-- name: FUnicodeToupper :one
SELECT Unicode::Toupper("hello");
-- name: FUnicodeReverse :one
SELECT Unicode::Reverse("hello");
-- name: FUnicodeIsascii :one
SELECT Unicode::Isascii("hello");
-- name: FUnicodeIsspace :one
SELECT Unicode::Isspace(" ");
-- name: FUnicodeIsupper :one
SELECT Unicode::Isupper("HELLO");
-- name: FUnicodeIslower :one
SELECT Unicode::Islower("hello");
-- name: FUnicodeIsalpha :one
SELECT Unicode::Isalpha("hello");
-- name: FUnicodeIsalnum :one
SELECT Unicode::Isalnum("hello123");
-- name: FUnicodeIshex :one
SELECT Unicode::Ishex("FF");
-- name: FUnicodeTouint64 :one
SELECT Unicode::Touint64("123");
-- name: FUnicodeLevensteindistance :one
SELECT Unicode::Levensteindistance("hello", "hallo");
