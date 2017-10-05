// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package controllers

import (
	"strconv"

	"github.com/EGaaS/go-egaas-mvp/packages/converter"
	"github.com/EGaaS/go-egaas-mvp/packages/model"
	"github.com/EGaaS/go-egaas-mvp/packages/utils"
)

// TxStatus return the status of the transaction
func (c *Controller) TxStatus() (string, error) {
	hash := c.r.FormValue("hash")
	ts := &model.TransactionStatus{}
	_, err := ts.Get(converter.HexToBin([]byte(hash)))
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if ts.BlockID != 0 {
		return `{"success":"` + strconv.FormatInt(ts.BlockID, 10) + `"}`, nil
	} else if len(ts.Error) > 0 {
		return "", utils.ErrInfo(ts.Error)
	}
	return `{"wait":"1"}`, nil
}