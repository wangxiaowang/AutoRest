package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

$typeObjectstruct$

func Get$Object$(c *gin.Context) {
	var (
		$object$ $Object$
		result gin.H
	)
	id := c.Param("id")
	row := db.QueryRow("select $TABLE_ELEMENTS$ from $TABLE$ where id = ?;", id)
	err := row.Scan($OBJECT_ELEMENTS$)
	if err != nil {
		// If no results send null
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": $object$,
			"count":  1,
		}
	}
	c.JSON(http.StatusOK, result)
}
func Get$Object$s(c *gin.Context) {
	var (
		$object$  $Object$
		$object$s []$Object$
	)
	rows, err := db.Query("select $TABLE_ELEMENTS$ from $TABLE$;")
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan($OBJECT_ELEMENTS$)
		$object$s = append($object$s, $object$)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": $object$s,
		"count":  len($object$s),
	})

}

func Post$Object$(c *gin.Context) {
	var buffer bytes.Buffer
	$POSTFORM_DATA$
	stmt, err := db.Prepare("insert into $TABLE$ ($POSTFORM_ELEMENTS$) values($POSTFORM_VALUE$);")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec($POSTFORM_ELEMENTS$)

	if err != nil {
		fmt.Print(err.Error())
	}

	$BUFFER_WRITE_STRING$
	defer stmt.Close()
	_name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf(" %s successfully created", _name),
	})
}
func Put$Object$(c *gin.Context) {
	var buffer bytes.Buffer
	id := c.Query("id")
	$POSTFORM_DATA$
	stmt, err := db.Prepare("update $TABLE$ set $PUT_WHERE$ where id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec($PUT_ELEMENTS$, id)
	if err != nil {
		fmt.Print(err.Error())
	}

	// Fastest way to append strings
	$BUFFER_WRITE_STRING$
	defer stmt.Close()
	_name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully updated to %s", _name),
	})

}
func Delete$Object$(c *gin.Context) {
	id := c.Query("id")
	stmt, err := db.Prepare("delete from $TABLE$ where id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Print(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted user: %s", id),
	})

}
