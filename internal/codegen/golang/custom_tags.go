package golang

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

const tagKeyValueSplit = "__@CUSTOM_TAG_VALUE__"

var tagReg = regexp.MustCompile(`^@.*\:`)

// Comment parse go struct tags
//
// Example:
// COMMENT ON COLUMN account.id is '@gorm:primaryKey @validate:required,min=3,max=32';
//
// To:
// type Account {
//   id int64 `gorm:"primaryKey" validate:"required,min=3,max=32"`
// }
//
func customTags(tags *map[string]string, column *catalog.Column) {
	if column.Comment == "" {
		return
	}

	comments := strings.Split(column.Comment, " ")
	fliterComments := make([]string, len(comments))

	for _, tag := range comments {
		if tagReg.Match([]byte(tag)) {
			tag = strings.Replace(tag, "@", "", 1)
			tag := strings.Replace(tag, ":", tagKeyValueSplit, 1)
			kv := strings.Split(tag, tagKeyValueSplit)
			if len(kv) < 2 {
				panic(fmt.Sprintf("comment tags JSON tags style error:  %s %s in %s", column.Type.Name, column.Name, column.Comment))
			}
			k := kv[0]
			v := kv[1]
			key := k + ":"
			(*tags)[key] = v
		} else {
			fliterComments = append(fliterComments, tag)
		}
	}

	// clear tags in Comment
	if len(fliterComments) > 0 {
		column.Comment = strings.Join(fliterComments, " ")
		column.Comment = strings.Trim(column.Comment, " ")
	}
}
