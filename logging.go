/*
 *  irc: IRC client library in Go
 *  Copyright (C) 2013  Joshua Chase <jcjoshuachase@gmail.com>
 *
 *  This program is free software; you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation; either version 2 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License along
 *  with this program; if not, write to the Free Software Foundation, Inc.,
 *  51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
*/

package irc

import (
	"github.com/Pursuit92/LeveledLogger/log"
	"io"
)

func SetLogLevel(n int) {
	log.Out.SetLevel(n)
}

func SetOut(wr io.Writer) {
	log.Out.SetOutput(wr)
}

func SetErr(wr io.Writer) {
	log.Err.SetOutput(wr)
}
