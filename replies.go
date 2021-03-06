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
	RplWelcome         string = "001"
	RplYourhost        string = "002"
	RplCreated         string = "003"
	RplMyinfo          string = "004"
	RplBounce          string = "005"
	RplUserhost        string = "302"
	RplIson            string = "303"
	RplAway            string = "301"
	RplUnaway          string = "305"
	RplNowaway         string = "306"
	RplWhoisuser       string = "311"
	RplWhoisserver     string = "312"
	RplWhoisoperator   string = "313"
	RplWhoisidle       string = "317"
	RplEndofwhois      string = "318"
	RplWhoischannels   string = "319"
	RplWhowasuser      string = "314"
	RplEndofwhowas     string = "369"
	RplListstart       string = "321"
	RplList            string = "322"
	RplListend         string = "323"
	RplUniqopis        string = "325"
	RplChannelmodeis   string = "324"
	RplNotopic         string = "331"
	RplTopic           string = "332"
	RplInviting        string = "341"
	RplSummoning       string = "342"
	RplInvitelist      string = "346"
	RplEndofinvitelist string = "347"
	RplExceptlist      string = "348"
	RplEndofexceptlist string = "349"
	RplVersion         string = "351"
	RplWhoreply        string = "352"
	RplEndofwho        string = "315"
	RplNamreply        string = "353"
	RplEndofnames      string = "366"
	RplLinks           string = "364"
	RplEndoflinks      string = "365"
	RplBanlist         string = "367"
	RplEndofbanlist    string = "368"
	RplInfo            string = "371"
	RplEndofinfo       string = "374"
	RplMotdstart       string = "375"
	RplMotd            string = "372"
	RplEndofmotd       string = "376"
	RplYoureoper       string = "381"
	RplRehashing       string = "382"
	RplYoureservice    string = "383"
	RplTime            string = "391"
	RplUserstart       string = "392"
	RplUsers           string = "393"
	RplNousers         string = "395"
	RplTracelink       string = "200"
	RplTraceconnecting string = "201"
	RplTracehandshake  string = "202"
	RplTraceunknown    string = "203"
	RplTraceoperator   string = "204"
	RplTraceuser       string = "205"
	RplTraceserver     string = "206"
	RplTraceservice    string = "207"
	RplTracenewtype    string = "208"
	RplTraceclass      string = "209"
	RplTracereconnect  string = "210"
	RplTracelog        string = "261"
	RplTraceend        string = "262"
	RplStatslinkinfo   string = "211"
	RplStatscommands   string = "212"
	RplEndofstats      string = "219"
	RplStatsuptime     string = "242"
	RplStatsoline      string = "243"
	RplUmodeis         string = "221"
	RplServlist        string = "234"
	RplServlistend     string = "235"
	RplLuserclient     string = "251"
	RplLuserop         string = "252"
	RplLuserunknown    string = "253"
	RplLuserchannels   string = "254"
	RplLuserme         string = "255"
	RplAdminme         string = "256"
	RplAdminloc1       string = "257"
	RplAdminloc2       string = "258"
	RplAdminemail      string = "259"
	RplTryagain        string = "263"
)
