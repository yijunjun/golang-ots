/*go**************************************************************************
 File            : main_test.go
 Subsystem       :
 Author          : yijunjun
 Date&Time       : 2017-03-01
 Description     :
 Revision        :

 History
 -------


 Copyright (c) Shenzhen Team Blemobi.
**************************************************************************go*/
package main

import (
	"testing"
)

func testSave(t *testing.T) {
	err := Save(&TFile{
		Name:       "test.xlsx",
		Sheets:     []string{"第一页", "第二页", "aaa", "1222"},
		RowHeaders: []string{"第一列头", "第三列头", "col", "343"},
		Rows: [][]string{
			[]string{"11", "22", "33", "44"},
			[]string{"aa", "bb", "cc", "dd"},
			[]string{"a1", "b2", "c3", "d4"},
		},
	})
	if err != nil {
		t.Error(err)
	}
}
