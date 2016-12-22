package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/eleme/esm-agent/log"
)

type Server struct {
	cfg *Config
}

var Prj string

func NewServer(c *Config) *Server {
	srv := &Server{cfg: c}
	return srv
}
func (s *Server) Open() error {
	var path string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}
	dir, _ := os.Getwd() //当前的目录
	Prj = dir + path + "project"
	err := os.Mkdir(Prj, os.ModePerm) //在当前目录下生成md目录
	if err != nil {
		log.Println(err)
		return err
	}

	CopyFile("tpl/MAIN.tpl", Prj+path+"main.go")
	for _, v := range s.cfg.JsonCfg {
		var table map[string]string
		if err := json.Unmarshal([]byte(v), &table); err != nil {
			return err
		}
		if name, ok := table["table_name"]; ok {
			CopyFile("tpl/OBJECT.tpl", Prj+path+name+".go")
		}
	}

	return nil

}
func (s *Server) Close() error {
	/*
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Println("delet dir error:", err)
		}*/
	return nil
}
func CopyFile(src, dst string) (w int64, rr error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println(err.Error())
		return 0, err
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}
func FirstUpper(name string) string {
	object := strings.ToLower(name)
	old := fmt.Sprintf("%c", object[0])
	new := fmt.Sprintf("%s", string(int(object[0])-32))
	Object := strings.Replace(object, old, new, 1)
	return Object
}

func (s *Server) TranslateTpl() error {
	//替换main.go中的模版
	f, err := ioutil.ReadFile("project/main.go")
	if err != nil {
		log.Println(err)
		return err
	}
	var main_code string
	main_code = string(f)
	main_code = strings.Replace(main_code, "$USER$", s.cfg.User, -1)
	main_code = strings.Replace(main_code, "$PASSWORD$", s.cfg.Password, -1)
	main_code = strings.Replace(main_code, "$HOST$", s.cfg.Host, -1)
	main_code = strings.Replace(main_code, "$PORT$", s.cfg.Port, -1)
	main_code = strings.Replace(main_code, "$DB$", s.cfg.DB, -1)

	var routerTpl = `router.GET("/$object$/:id", func(c *gin.Context) {
		Get$Object$(c)
	})
	router.GET("/$object$s", func(c *gin.Context) {
		Get$Object$s(c)
	})
	router.POST("/$object$", func(c *gin.Context) {
		Post$Object$(c)
	})
	router.PUT("/$object$", func(c *gin.Context) {
		Put$Object$(c)
	})
	router.DELETE("/$object$", func(c *gin.Context) {
		Delete$Object$(c)
	})`

	var TotalRouterStr string
	for _, v := range s.cfg.JsonCfg {
		var table map[string]string
		if err := json.Unmarshal([]byte(v), &table); err != nil {
			return err
		}
		if name, ok := table["table_name"]; ok {
			object := strings.ToLower(name)
			Object := FirstUpper(object)
			routerStr := strings.Replace(routerTpl, "$object$", object, -1)
			routerStr = strings.Replace(routerStr, "$Object$", Object, -1)
			TotalRouterStr += routerStr + "\n"
		}
	}
	main_code = strings.Replace(main_code, "$ROUTERS$", TotalRouterStr, -1)
	err = ioutil.WriteFile("project/main.go", []byte(main_code), 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//替换每个对象中的模板
	for _, v := range s.cfg.JsonCfg {
		var table map[string]string
		if err := json.Unmarshal([]byte(v), &table); err != nil {
			return err
		}
		if name, ok := table["table_name"]; ok {
			object := strings.ToLower(name)
			f, err := ioutil.ReadFile("project/" + object + ".go")
			if err != nil {
				log.Println(err)
				return err
			}
			totalStr := string(f)
			Object := FirstUpper(object)
			totalStr = strings.Replace(totalStr, "$object$", object, -1)
			totalStr = strings.Replace(totalStr, "$Object$", Object, -1)

			var elements = make([]string, 0)
			delete(table, "table_name")
			for k, _ := range table {

				elements = append(elements, k)
			}
			sort.Strings(elements)
			//$TABLE$
			totalStr = strings.Replace(totalStr, "$TABLE$", object, -1)
			//$typeObjectstruct$
			typeObjectstruct := "type " + Object + " struct{\n"
			typeObjectstruct += "Id  int\n"

			for _, v := range elements {
				typeObjectstruct += FirstUpper(v) + "   " + table[v] + "\n"
			}
			typeObjectstruct += "}"

			totalStr = strings.Replace(totalStr, "$typeObjectstruct$", typeObjectstruct, -1)

			//$TABLE_ELEMENTS$
			table_elements := "id"
			for _, v := range elements {
				table_elements += ","
				table_elements += strings.ToLower(v)
			}

			totalStr = strings.Replace(totalStr, "$TABLE_ELEMENTS$", table_elements, -1)

			//$OBJECT_ELEMENTS$
			Object_elements := "&" + object + ".Id,"
			for k, v := range elements {
				Object_elements += "&"
				Object_elements += object
				Object_elements += "."
				Object_elements += FirstUpper(v)
				if k != len(elements)-1 {
					Object_elements += ","
				}
			}
			totalStr = strings.Replace(totalStr, "$OBJECT_ELEMENTS$", Object_elements, -1)

			//$POSTFORM_DATA$
			postForm_data := ""
			for _, v := range elements {
				postForm_data += strings.ToLower(v) + ":=c.PostForm(\"" + strings.ToLower(v) + "\")\n"
			}

			totalStr = strings.Replace(totalStr, "$POSTFORM_DATA$", postForm_data, -1)

			//$POSTFORM_ELEMENTS$
			postForm_elements := ""
			for k, v := range elements {
				postForm_elements += strings.ToLower(v)
				if k != len(elements)-1 {
					postForm_elements += ","
				}
			}

			totalStr = strings.Replace(totalStr, "$POSTFORM_ELEMENTS$", postForm_elements, -1)
			//$POSTFORM_VALUE$
			postForm_value := ""
			for k, _ := range elements {
				postForm_value += "?"
				if k != len(elements)-1 {
					postForm_value += ","
				}
			}

			totalStr = strings.Replace(totalStr, "$POSTFORM_VALUE$", postForm_value, -1)
			//$BUFFER_WRITE_STRING$
			buffer_write_string := ""
			for _, v := range elements {
				buffer_write_string += "buffer.WriteString(" + strings.ToLower(v) + ")\nbuffer.WriteString(\"  \")\n"
			}

			totalStr = strings.Replace(totalStr, "$BUFFER_WRITE_STRING$", buffer_write_string, -1)

			//$PUT_WHERE$
			put_where := ""
			for k, v := range elements {
				put_where += strings.ToLower(v) + "=?"
				if k != len(elements)-1 {
					put_where += ","
				}
			}

			totalStr = strings.Replace(totalStr, "$PUT_WHERE$", put_where, -1)
			//$PUT_ELEMENTS$
			put_elements := postForm_elements
			totalStr = strings.Replace(totalStr, "$PUT_ELEMENTS$", put_elements, -1)
			err = ioutil.WriteFile("project/"+object+".go", []byte(totalStr), 0666)
			if err != nil {
				fmt.Println(err)
				return err
			}

		}
	}
	return nil

}

func (s *Server) ToSQL() error {
	//通过配置的对象生成SQL文件
	/*
			SQLTpl:="CREATE DATABASE $DB$ IF NOT EXISTS `$DB$`;
			      use $DB$;
			    DROP TABLE IF EXISTS `$TABLE$`;
				   CREATE TABLE `collector_types` (
				      `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
			       `name` varchar(255) NOT NULL,
			    `is_enabled` tinyint(1) NOT NULL,
		     `create_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		      `update_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		       PRIMARY KEY (`id`),
		    UNIQUE KEY `uix_collector_types_name` (`name`),
		     KEY `idx_collector_types_is_enabled` (`is_enabled`)
		    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;"

			for _, v := range s.cfg.JsonCfg {
				var table map[string]string
				if err := json.Unmarshal([]byte(v), &table); err != nil {
					return err
				}
				if name, ok := table["table_name"]; ok {
					object := strings.ToLower(name)

				}
			}*/
	return nil
}
func (s *Server) GoBuild() {
	//编译代码生成二进制文件

}
