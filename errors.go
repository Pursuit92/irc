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

const (
	ErrNosuchnick        string = "401"
	ErrNosuchserver      string = "402"
	ErrNosuchchannel     string = "403"
	ErrCannotsendtochan  string = "404"
	ErrToomanychannels   string = "405"
	ErrWasnosuchnick     string = "406"
	ErrToomanytargets    string = "407"
	ErrNosuchservice     string = "408"
	ErrNoorigin          string = "409"
	ErrNorecipient       string = "411"
	ErrNotexttosend      string = "412"
	ErrNotoplevel        string = "413"
	ErrWildtoplevel      string = "414"
	ErrBadmask           string = "415"
	ErrUnknowncommand    string = "421"
	ErrNomotd            string = "422"
	ErrNoadmininfo       string = "423"
	ErrFileerror         string = "424"
	ErrNonicknamegiven   string = "431"
	ErrErroneusnickname  string = "432"
	ErrNicknameinuse     string = "433"
	ErrNickcollision     string = "436"
	ErrUnavailresource   string = "437"
	ErrUsernotinchannel  string = "441"
	ErrNotonchannel      string = "442"
	ErrUseronchannel     string = "443"
	ErrNologin           string = "444"
	ErrSummondisabled    string = "445"
	ErrUserdisabled      string = "446"
	ErrNotregistered     string = "451"
	ErrNeedmoreparams    string = "461"
	ErrAlreadyRegistered string = "462"
	ErrNopermforhost     string = "463"
	ErrPasswdmismatch    string = "464"
	ErrYourbannedcreep   string = "465"
	ErrYouwillbebanned   string = "466"
	ErrKeyset            string = "467"
	ErrChannelisfull     string = "471"
	ErrUnknownmode       string = "472"
	ErrInviteonlychan    string = "473"
	ErrBannedfromchan    string = "474"
	ErrBadchannelkey     string = "475"
	ErrBadchanmask       string = "476"
	ErrNochanmodes       string = "477"
	ErrBanlistfull       string = "478"
	ErrNoprivileges      string = "481"
	ErrChanoprivsneeded  string = "482"
	ErrCantkillserver    string = "483"
	ErrRestricted        string = "484"
	ErrUniqoprivsneeded  string = "485"
	ErrNooperhost        string = "491"
	ErrUmodeunknownflag  string = "501"
	ErrUsersdontmatch    string = "502"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	UserTaken = Error("Username already taken.")
	InvalidPass = Error("Invalid password.")
	Timeout = Error("Connect timeout.")
	Disconnect = Error("Remote disconnected.")
)
