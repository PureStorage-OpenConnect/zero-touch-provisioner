/*
	Pure Storage FlashArray and FlashBlade Zero Touch Provisioner
	Brandon Showers
	March 22 2020
	v1
*/

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/buger/jsonparser"
	"gopkg.in/go-playground/validator.v9"
)

//FLASH ARRAY VARS//
var mainwin *ui.Window
var ipAddress = ""

//END FLASH ARRAY VARS//

//FLASH BLADE VARS//
//Global Vars//
var ipAddressFB = ""
var xAuthToken = ""
var loginUrl = ""
var apiUrl = ""
var steps = make(map[int]string)
var progressCounter = 0
var statusCode int

//END FLASH BLADE VARS//

//FUNCTIONS//
func timeZones() []string {
	tz := []string{"Africa/Abidjan", "Africa/Accra", "Africa/Addis_Ababa", "Africa/Algiers", "Africa/Asmara", "Africa/Bamako", "Africa/Bangui", "Africa/Banjul", "Africa/Bissau", "Africa/Blantyre", "Africa/Brazzaville", "Africa/Bujumbura", "Africa/Cairo", "Africa/Casablanca", "Africa/Ceuta", "Africa/Conakry", "Africa/Dakar", "Africa/Dar_es_Salaam", "Africa/Djibouti", "Africa/Douala", "Africa/El_Aaiun", "Africa/Freetown", "Africa/Gaborone", "Africa/Harare", "Africa/Johannesburg", "Africa/Juba", "Africa/Kampala", "Africa/Khartoum", "Africa/Kigali", "Africa/Kinshasa", "Africa/Lagos", "Africa/Libreville", "Africa/Lome", "Africa/Luanda", "Africa/Lubumbashi", "Africa/Lusaka", "Africa/Malabo", "Africa/Maputo", "Africa/Maseru", "Africa/Mbabane", "Africa/Mogadishu", "Africa/Monrovia", "Africa/Nairobi", "Africa/Ndjamena", "Africa/Niamey", "Africa/Nouakchott", "Africa/Ouagadougou", "Africa/Porto-Novo", "Africa/Sao_Tome", "Africa/Tripoli", "Africa/Tunis", "Africa/Windhoek", "America/Adak", "America/Anchorage", "America/Anguilla", "America/Antigua", "America/Araguaina", "America/Argentina/Buenos_Aires", "America/Argentina/Catamarca", "America/Argentina/Cordoba", "America/Argentina/Jujuy", "America/Argentina/La_Rioja", "America/Argentina/Mendoza", "America/Argentina/Rio_Gallegos", "America/Argentina/Salta", "America/Argentina/San_Juan", "America/Argentina/San_Luis", "America/Argentina/Tucuman", "America/Argentina/Ushuaia", "America/Aruba", "America/Asuncion", "America/Atikokan", "America/Bahia", "America/Bahia_Banderas", "America/Barbados", "America/Belem", "America/Belize", "America/Blanc-Sablon", "America/Boa_Vista", "America/Bogota", "America/Boise", "America/Cambridge_Bay", "America/Campo_Grande", "America/Cancun", "America/Caracas", "America/Cayenne", "America/Cayman", "America/Chicago", "America/Chihuahua", "America/Costa_Rica", "America/Creston", "America/Cuiaba", "America/Curacao", "America/Danmarkshavn", "America/Dawson", "America/Dawson_Creek", "America/Denver", "America/Detroit", "America/Dominica", "America/Edmonton", "America/Eirunepe", "America/El_Salvador", "America/Fort_Nelson", "America/Fortaleza", "America/Glace_Bay", "America/Godthab", "America/Goose_Bay", "America/Grand_Turk", "America/Grenada", "America/Guadeloupe", "America/Guatemala", "America/Guayaquil", "America/Guyana", "America/Halifax", "America/Havana", "America/Hermosillo", "America/Indiana/Indianapolis", "America/Indiana/Knox", "America/Indiana/Marengo", "America/Indiana/Petersburg", "America/Indiana/Tell_City", "America/Indiana/Vevay", "America/Indiana/Vincennes", "America/Indiana/Winamac", "America/Inuvik", "America/Iqaluit", "America/Jamaica", "America/Juneau", "America/Kentucky/Louisville", "America/Kentucky/Monticello", "America/Kralendijk", "America/La_Paz", "America/Lima", "America/Los_Angeles", "America/Lower_Princes", "America/Maceio", "America/Managua", "America/Manaus", "America/Marigot", "America/Martinique", "America/Matamoros", "America/Mazatlan", "America/Menominee", "America/Merida", "America/Metlakatla", "America/Mexico_City", "America/Miquelon", "America/Moncton", "America/Monterrey", "America/Montevideo", "America/Montserrat", "America/Nassau", "America/New_York", "America/Nipigon", "America/Nome", "America/Noronha", "America/North_Dakota/Beulah", "America/North_Dakota/Center", "America/North_Dakota/New_Salem", "America/Ojinaga", "America/Panama", "America/Pangnirtung", "America/Paramaribo", "America/Phoenix", "America/Port_of_Spain", "America/Port-au-Prince", "America/Porto_Velho", "America/Puerto_Rico", "America/Punta_Arenas", "America/Rainy_River", "America/Rankin_Inlet", "America/Recife", "America/Regina", "America/Resolute", "America/Rio_Branco", "America/Santarem", "America/Santiago", "America/Santo_Domingo", "America/Sao_Paulo", "America/Scoresbysund", "America/Sitka", "America/St_Barthelemy", "America/St_Johns", "America/St_Kitts", "America/St_Lucia", "America/St_Thomas", "America/St_Vincent", "America/Swift_Current", "America/Tegucigalpa", "America/Thule", "America/Thunder_Bay", "America/Tijuana", "America/Toronto", "America/Tortola", "America/Vancouver", "America/Whitehorse", "America/Winnipeg", "America/Yakutat", "America/Yellowknife", "Antarctica/Casey", "Antarctica/Davis", "Antarctica/DumontDUrville", "Antarctica/Macquarie", "Antarctica/Mawson", "Antarctica/McMurdo", "Antarctica/Palmer", "Antarctica/Rothera", "Antarctica/Syowa", "Antarctica/Troll", "Antarctica/Vostok", "Arctic/Longyearbyen", "Asia/Aden", "Asia/Almaty", "Asia/Amman", "Asia/Anadyr", "Asia/Aqtau", "Asia/Aqtobe", "Asia/Ashgabat", "Asia/Atyrau", "Asia/Baghdad", "Asia/Bahrain", "Asia/Baku", "Asia/Bangkok", "Asia/Barnaul", "Asia/Beirut", "Asia/Bishkek", "Asia/Brunei", "Asia/Chita", "Asia/Choibalsan", "Asia/Colombo", "Asia/Damascus", "Asia/Dhaka", "Asia/Dili", "Asia/Dubai", "Asia/Dushanbe", "Asia/Famagusta", "Asia/Gaza", "Asia/Hebron", "Asia/Ho_Chi_Minh", "Asia/Hong_Kong", "Asia/Hovd", "Asia/Irkutsk", "Asia/Jakarta", "Asia/Jayapura", "Asia/Jerusalem", "Asia/Kabul", "Asia/Kamchatka", "Asia/Karachi", "Asia/Kathmandu", "Asia/Khandyga", "Asia/Kolkata", "Asia/Krasnoyarsk", "Asia/Kuala_Lumpur", "Asia/Kuching", "Asia/Kuwait", "Asia/Macau", "Asia/Magadan", "Asia/Makassar", "Asia/Manila", "Asia/Muscat", "Asia/Nicosia", "Asia/Novokuznetsk", "Asia/Novosibirsk", "Asia/Omsk", "Asia/Oral", "Asia/Phnom_Penh", "Asia/Pontianak", "Asia/Pyongyang", "Asia/Qatar", "Asia/Qostanay", "Asia/Qyzylorda", "Asia/Riyadh", "Asia/Sakhalin", "Asia/Samarkand", "Asia/Seoul", "Asia/Shanghai", "Asia/Singapore", "Asia/Srednekolymsk", "Asia/Taipei", "Asia/Tashkent", "Asia/Tbilisi", "Asia/Tehran", "Asia/Thimphu", "Asia/Tokyo", "Asia/Tomsk", "Asia/Ulaanbaatar", "Asia/Urumqi", "Asia/Ust-Nera", "Asia/Vientiane", "Asia/Vladivostok", "Asia/Yakutsk", "Asia/Yangon", "Asia/Yekaterinburg", "Asia/Yerevan", "Atlantic/Azores", "Atlantic/Bermuda", "Atlantic/Canary", "Atlantic/Cape_Verde", "Atlantic/Faroe", "Atlantic/Madeira", "Atlantic/Reykjavik", "Atlantic/South_Georgia", "Atlantic/St_Helena", "Atlantic/Stanley", "Australia/Adelaide", "Australia/Brisbane", "Australia/Broken_Hill", "Australia/Currie", "Australia/Darwin", "Australia/Eucla", "Australia/Hobart", "Australia/Lindeman", "Australia/Lord_Howe", "Australia/Melbourne", "Australia/Perth", "Australia/Sydney", "Europe/Amsterdam", "Europe/Andorra", "Europe/Astrakhan", "Europe/Athens", "Europe/Belgrade", "Europe/Berlin", "Europe/Bratislava", "Europe/Brussels", "Europe/Bucharest", "Europe/Budapest", "Europe/Busingen", "Europe/Chisinau", "Europe/Copenhagen", "Europe/Dublin", "Europe/Gibraltar", "Europe/Guernsey", "Europe/Helsinki", "Europe/Isle_of_Man", "Europe/Istanbul", "Europe/Jersey", "Europe/Kaliningrad", "Europe/Kiev", "Europe/Kirov", "Europe/Lisbon", "Europe/Ljubljana", "Europe/London", "Europe/Luxembourg", "Europe/Madrid", "Europe/Malta", "Europe/Mariehamn", "Europe/Minsk", "Europe/Monaco", "Europe/Moscow", "Europe/Oslo", "Europe/Paris", "Europe/Podgorica", "Europe/Prague", "Europe/Riga", "Europe/Rome", "Europe/Samara", "Europe/San_Marino", "Europe/Sarajevo", "Europe/Saratov", "Europe/Simferopol", "Europe/Skopje", "Europe/Sofia", "Europe/Stockholm", "Europe/Tallinn", "Europe/Tirane", "Europe/Ulyanovsk", "Europe/Uzhgorod", "Europe/Vaduz", "Europe/Vatican", "Europe/Vienna", "Europe/Vilnius", "Europe/Volgograd", "Europe/Warsaw", "Europe/Zagreb", "Europe/Zaporozhye", "Europe/Zurich", "Indian/Antananarivo", "Indian/Chagos", "Indian/Christmas", "Indian/Cocos", "Indian/Comoro", "Indian/Kerguelen", "Indian/Mahe", "Indian/Maldives", "Indian/Mauritius", "Indian/Mayotte", "Indian/Reunion", "Pacific/Apia", "Pacific/Auckland", "Pacific/Bougainville", "Pacific/Chatham", "Pacific/Chuuk", "Pacific/Easter", "Pacific/Efate", "Pacific/Enderbury", "Pacific/Fakaofo", "Pacific/Fiji", "Pacific/Funafuti", "Pacific/Galapagos", "Pacific/Gambier", "Pacific/Guadalcanal", "Pacific/Guam", "Pacific/Honolulu", "Pacific/Kiritimati", "Pacific/Kosrae", "Pacific/Kwajalein", "Pacific/Majuro", "Pacific/Marquesas", "Pacific/Midway", "Pacific/Nauru", "Pacific/Niue", "Pacific/Norfolk", "Pacific/Noumea", "Pacific/Pago_Pago", "Pacific/Palau", "Pacific/Pitcairn", "Pacific/Pohnpei", "Pacific/Port_Moresby", "Pacific/Rarotonga", "Pacific/Saipan", "Pacific/Tahiti", "Pacific/Tarawa", "Pacific/Tongatapu", "Pacific/Wake", "Pacific/Wallis"}
	return tz
}

//Post rest function specifically for FB logon only takes 2 parameters and returns a string//
func postAPICallLoginFB(url string, apiToken string) string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	req.Header.Set("api-token", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}

	//set the status code for the response
	statusCode = resp.StatusCode

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}

	//Sets the x-auth-token from the header response
	if len(resp.Header["X-Auth-Token"]) > 0 {
		s := resp.Header["X-Auth-Token"]
		t := strings.Replace(s[0], "[", "", -1)
		t = strings.Replace(t, "]", "", -1)
		xAuthToken = t
	}

	return string(body)
}

//PRIMARY REST METHOD//
func apiCallFB(method, url string, xAuthToken string, data []byte) []byte {
	//data is only used for post and patch. for delete and get this is nil
	jsonBody := bytes.NewReader(data)
	//new http client that ignores ssl cert errors.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
	req, err := http.NewRequest(method, url, jsonBody)
	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-auth-token", xAuthToken)

	//make the rest call
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}
	//wait then close the connection to free space.
	defer resp.Body.Close()

	//set the status code for the response
	statusCode = resp.StatusCode

	//convert http.response body to byte array
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println(err2)
		return []byte(err2.Error())
	}
	//finally return the response byte array.
	return body
}

//END FUNCTIONS//

//MAIN FLASHARRAY UI WORKER FUNCTION//
func initializeFATab() ui.Control {
	//fields for the form
	arrayName := ui.NewEntry()
	eulaOrg := ui.NewEntry()
	eulaName := ui.NewEntry()
	eulaTitle := ui.NewEntry()
	eulaAccept := ui.NewCheckbox("yes")
	ntpServer := ui.NewEntry()
	vir0IP := ui.NewEntry()
	vir0SNM := ui.NewEntry()
	vir0GW := ui.NewEntry()
	ct0IP := ui.NewEntry()
	ct0SNM := ui.NewEntry()
	ct0GW := ui.NewEntry()
	ct1IP := ui.NewEntry()
	ct1SNM := ui.NewEntry()
	ct1GW := ui.NewEntry()
	dnsDomain := ui.NewEntry()
	dnsServer := ui.NewEntry()
	smtpRelay := ui.NewEntry()
	smtpDomain := ui.NewEntry()
	smtpAlertEmail := ui.NewEntry()
	tempIP := ui.NewEntry() //dhcp ip address
	initResult := ui.NewMultilineEntry()
	timeZone := ui.NewCombobox()
	tz := timeZones()
	for i, v := range tz {
		timeZone.Append(v)
		i++
	}

	//first column definition
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	//define vertical box inside column similar to a div
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	//ARRAY NAME FIELD//
	//define the group for the form
	group1 := ui.NewGroup("General Configs")
	group1.SetMargined(true)

	//add group to the vertical box
	vbox.Append(group1, false)

	//define the form for the group
	entryForm1 := ui.NewForm()
	entryForm1.SetPadded(true)

	//embed the array name form field inside the first form group
	group1.SetChild(entryForm1)
	entryForm1.Append("FlashArray Name", arrayName, false)
	entryForm1.Append("", ui.NewLabel(""), false)
	entryForm1.Append("Organization Name", eulaOrg, false)
	entryForm1.Append("Your Name", eulaName, false)
	entryForm1.Append("Your Title", eulaTitle, false)
	entryForm1.Append("You accept EULA", eulaAccept, false)
	entryForm1.Append("", ui.NewLabel(""), false)
	entryForm1.Append("NTP Time Server(s)**", ntpServer, false)
	entryForm1.Append("TimeZone", timeZone, false)
	entryForm1.Append("", ui.NewLabel("  ________Optional Below________ "), false)
	entryForm1.Append("", ui.NewLabel(""), false)
	entryForm1.Append("DNS Domain", dnsDomain, false)
	entryForm1.Append("DNS Name Server(s)**", dnsServer, false)
	entryForm1.Append("", ui.NewLabel(""), false)
	entryForm1.Append("SMTP Relay Host", smtpRelay, false)
	entryForm1.Append("SMTP sender domain", smtpDomain, false)
	entryForm1.Append("Alert Email Address(s)**", smtpAlertEmail, false)
	entryForm1.Append("", ui.NewLabel("**Comma seperated"), false)

	//seperator line
	hbox.Append(ui.NewVerticalSeparator(), false)

	//Middle column
	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	//VIR0IP FORM//
	group3 := ui.NewGroup("Virtual Nic 0")
	group3.SetMargined(true)
	vbox.Append(group3, false)

	entryForm3 := ui.NewForm()
	entryForm3.SetPadded(true)
	group3.SetChild(entryForm3)

	//autofill button to copy contents to ct0 and ct1 ip configs
	button := ui.NewButton("Autofill")
	entryForm3.Append("IP Address", vir0IP, false)
	entryForm3.Append("Subnet Mask", vir0SNM, false)
	entryForm3.Append("Default Gateway", vir0GW, false)
	entryForm3.Append("Replicate below", button, false)

	//CT0 FORM//
	group5 := ui.NewGroup("Controller 0")
	group5.SetMargined(true)
	vbox.Append(group5, false)
	entryForm5 := ui.NewForm()
	entryForm5.SetPadded(true)
	group5.SetChild(entryForm5)

	entryForm5.Append("IP Address", ct0IP, false)
	entryForm5.Append("Subnet Mask", ct0SNM, false)
	entryForm5.Append("Default Gateway", ct0GW, false)

	//CT1 FORM//
	group6 := ui.NewGroup("Controller 1")
	group6.SetMargined(true)
	vbox.Append(group6, false)
	entryForm6 := ui.NewForm()
	entryForm6.SetPadded(true)
	group6.SetChild(entryForm6)

	entryForm6.Append("IP Address", ct1IP, false)
	entryForm6.Append("Subnet Mask", ct1SNM, false)
	entryForm6.Append("Default Gateway", ct1GW, false)

	hbox.Append(ui.NewVerticalSeparator(), false)

	//third column
	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, true)

	//SUBMIT "GO" BUTTON//
	group9 := ui.NewGroup("Initialize Array")
	group9.SetMargined(true)
	vbox.Append(group9, true)

	entryForm9 := ui.NewForm()
	entryForm9.SetPadded(true)
	group9.SetChild(entryForm9)

	button1 := ui.NewButton("Query")
	entryForm9.Append("", ui.NewLabel(""), false)

	//submit and go button
	button2 := ui.NewButton("Initialize")

	entryForm9.Append("DHCP IP of Array ", tempIP, false)
	entryForm9.Append("Query First, ", button1, false)
	entryForm9.Append("Configure Array ", button2, false)

	//multiline field for showing results of patch api call and form validation messages.
	//sets the initResults console to readonly
	initResult.SetReadOnly(true)
	//sets initial instructions into the console window.
	initResult.SetText("Welcome to the FlashArray Zero Touch Provisioner!\n\nYou should have obtained the DHCP IP of the recently installed FlashArray you will be initializing with this tool.  Enter it above and press the Query button to confirm your connectivity.\n\nWhen you are ready, fill out the form and press the Initialize button to configure your Array.\n\nAfter the Array is initialized, you will not be able to re-connect again with this tool.  You will need to use the CLI or GUI for additonal configuration.\n\nPlease contact Pure Support or your Account team for any questions or issues.")
	entryForm9.Append("Init Results", initResult, true)

	//autofill IP config button actions
	//used to replicate the ip info from vi0 to ct0 and ct1
	button.OnClicked(func(*ui.Button) {

		ct0IP.SetText(vir0IP.Text())
		ct0SNM.SetText(vir0SNM.Text())
		ct0GW.SetText(vir0GW.Text())
		ct1IP.SetText(vir0IP.Text())
		ct1SNM.SetText(vir0SNM.Text())
		ct1GW.SetText(vir0GW.Text())

	})

	button1.OnClicked(func(*ui.Button) {
		ipAddress = tempIP.Text()
		//query the FA
		result := apiCallFB("GET", "http://"+ipAddress+":8081/array-initial-config", "", nil)

		//testing only
		//result := apiCallFB("GET", "https://pureapisim.azurewebsites.net/api/array-initial-config", "", nil)

		//set results from apiCall to the initResult field.
		initResult.SetText(string(result))

	})

	//initialize the array and do lots of other work
	button2.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate Controller 1 Gateway
		err7 := validate.Var(ct1GW.Text(), "required,ipv4")
		if err7 != nil {
			initResult.SetText("Please provide a valid Gateway for Controller 1")
			passed = false
		}
		//validate Controller 1 SN
		err8 := validate.Var(ct1SNM.Text(), "required,ipv4")
		if err8 != nil {
			initResult.SetText("Please provide a valid Subnet Mask for Controller 1")
			passed = false
		}
		//validate Controller 1 IP
		err9 := validate.Var(ct1IP.Text(), "required,ipv4")
		if err9 != nil {
			initResult.SetText("Please provide a valid IP Address for Controller 1")
			passed = false
		}
		//validate Controller 0 Gateway
		err10 := validate.Var(ct0GW.Text(), "required,ipv4")
		if err10 != nil {
			initResult.SetText("Please provide a valid Gateway for Controller 0")
			passed = false
		}
		//validate Controller 0 SN
		err11 := validate.Var(ct0SNM.Text(), "required,ipv4")
		if err11 != nil {
			initResult.SetText("Please provide a valid Subnet Mask for Controller 0")
			passed = false
		}
		//validate Controller 0 IP
		err12 := validate.Var(ct0IP.Text(), "required,ipv4")
		if err12 != nil {
			initResult.SetText("Please provide a valid IP Address for Controller 0")
			passed = false
		}
		//validate Virtual 0 Gateway
		err13 := validate.Var(vir0GW.Text(), "required,ipv4")
		if err13 != nil {
			initResult.SetText("Please provide a valid Gateway for Virtual 0")
			passed = false
		}
		//validate Virtual 0 SN
		err14 := validate.Var(vir0SNM.Text(), "required,ipv4")
		if err14 != nil {
			initResult.SetText("Please provide a valid Subnet Mask for Virtual 0")
			passed = false
		}
		//validate Virtual 0 IP
		err15 := validate.Var(vir0IP.Text(), "required,ipv4")
		if err15 != nil {
			initResult.SetText("Please provide a valid IP Address for Virtual 0")
			passed = false
		}
		//validate SMTP Relay Host
		if smtpRelay.Text() != "" {
			err8 := validate.Var(smtpRelay.Text(), "fqdn|ip|url")
			if err8 != nil {
				initResult.SetText("Please a valid SMTP Relay Host using either FQDN,IP or URL.")
				passed = false
			}
		}
		//validate SMTP sender domain
		if smtpDomain.Text() != "" {
			err9 := validate.Var(smtpDomain.Text(), "fqdn")
			if err9 != nil {
				initResult.SetText("Please enter a FQDN for your SMTP sender domain.")
				passed = false
			}
		}
		//validate alert email addresses.
		ae := strings.Split(smtpAlertEmail.Text(), ",")
		if smtpAlertEmail.Text() != "" {
			for i := 0; i < len(ae); i++ {
				//fmt.Print(ntp[i] + "\n")
				err2 := validate.Var(ae[i], "email")
				if err2 != nil {
					initResult.SetText("Please provide a valid email address.\n\nIf more than one email address is entered please use comma seperation with no spaces in-between.")
					passed = false
				}
			}
		}
		//validate DNS servers if entered
		ns := strings.Split(dnsServer.Text(), ",")
		if dnsServer.Text() != "" {
			for i := 0; i < len(ns); i++ {
				err6 := validate.Var(ns[i], "fqdn|ip")
				if err6 != nil {
					initResult.SetText("Please provide a fqdn or ip for the DNS server.\n\nIf more than one server is entered please use comma seperation with no spaces in-between.")
					passed = false
				}
			}
		}
		//validate DNS Domain name
		if dnsDomain.Text() != "" {
			err5 := validate.Var(dnsDomain.Text(), "fqdn")
			if err5 != nil {
				initResult.SetText("Please a FQDN for your DNS Domain.")
				passed = false
			}
		}
		//validate Ntp server
		ntp := strings.Split(ntpServer.Text(), ",")
		for i := 0; i < len(ntp); i++ {
			//fmt.Print(ntp[i] + "\n")
			err7 := validate.Var(ntp[i], "fqdn|ip")
			if err7 != nil {
				initResult.SetText("Please provide a fqdn or ip for the NTP server.\n\nIf more than one server is entered please use comma seperation with no spaces in-between.")
				passed = false
			}
		}
		//validate timezone
		if timeZone.Selected() < 0 {
			initResult.SetText("Please select a Timezone")
			passed = false
		}
		//validate eula
		if eulaAccept.Checked() != true {
			initResult.SetText("You must accept the terms of our EULA")
			passed = false
		}
		//validate Eula Title
		err4 := validate.Var(eulaTitle.Text(), "required")
		if err4 != nil {
			initResult.SetText("Please provide your Job Title")
			passed = false
		}
		//validate Eula Name
		err3 := validate.Var(eulaName.Text(), "required")
		if err3 != nil {
			initResult.SetText("Please provide your Full Name")
			passed = false
		}
		//validate Eula Org Name
		err2 := validate.Var(eulaOrg.Text(), "required")
		if err2 != nil {
			initResult.SetText("Please provide your Organization Name")
			passed = false
		}
		//validate Array Name
		var rxPatArrayName = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,54}[a-zA-Z0-9])?$`)
		if !rxPatArrayName.MatchString(arrayName.Text()) {
			initResult.SetText("ArrayName has blank or contains invalid characters.  It must begin with a number or letter, can contain a dash in the body of the name, but must also end with a number or letter.   No more than 55 characters in length.")
		}
		//validate DHCP Boot IP
		err0 := validate.Var(tempIP.Text(), "required,ipv4")
		if err0 != nil {
			initResult.SetText("Please provide a valid IP Address for the DHCP boot IP")
			passed = false
		}
		//if all validation above passes then proceed...
		if passed == true {
			//cool site to generate struct from json https://mholt.github.io/json-to-go/
			//define the flash array json structure
			type FAS struct {
				ArrayName string `json:"array_name"`
				Ct0Eth0   struct {
					Address string `json:"address"`
					Netmask string `json:"netmask"`
					Gateway string `json:"gateway"`
				} `json:"ct0.eth0"`
				Ct1Eth0 struct {
					Address string `json:"address"`
					Netmask string `json:"netmask"`
					Gateway string `json:"gateway"`
				} `json:"ct1.eth0"`
				Vir0 struct {
					Address string `json:"address"`
					Netmask string `json:"netmask"`
					Gateway string `json:"gateway"`
				} `json:"vir0"`
				DNS struct {
					Domain      string   `json:"domain"`
					Nameservers []string `json:"nameservers"`
				} `json:"dns"`
				NtpServers []string `json:"ntp_servers"`
				Timezone   string   `json:"timezone"`
				SMTP       struct {
					RelayHost    string `json:"relay_host"`
					SenderDomain string `json:"sender_domain"`
				} `json:"smtp"`
				AlertEmails    []string `json:"alert_emails"`
				EulaAcceptance struct {
					Accepted   bool `json:"accepted"`
					AcceptedBy struct {
						Organization string `json:"organization"`
						FullName     string `json:"full_name"`
						JobTitle     string `json:"job_title"`
					} `json:"accepted_by"`
				} `json:"eula_acceptance"`
			}

			//initialize FAS struct object
			FA := &FAS{}
			FA.ArrayName = arrayName.Text()
			FA.Ct0Eth0.Address = ct0IP.Text()
			FA.Ct0Eth0.Netmask = ct0SNM.Text()
			FA.Ct0Eth0.Gateway = ct0GW.Text()
			FA.Ct1Eth0.Address = ct1IP.Text()
			FA.Ct1Eth0.Netmask = ct1SNM.Text()
			FA.Ct1Eth0.Gateway = ct1GW.Text()
			FA.Vir0.Address = vir0IP.Text()
			FA.Vir0.Netmask = vir0SNM.Text()
			FA.Vir0.Gateway = vir0GW.Text()
			FA.DNS.Domain = dnsDomain.Text()
			FA.DNS.Nameservers = ns
			FA.NtpServers = ntp
			FA.Timezone = tz[timeZone.Selected()]
			FA.SMTP.RelayHost = smtpRelay.Text()
			FA.SMTP.SenderDomain = smtpDomain.Text()
			FA.AlertEmails = ae
			FA.EulaAcceptance.Accepted = eulaAccept.Checked()
			FA.EulaAcceptance.AcceptedBy.FullName = eulaName.Text()
			FA.EulaAcceptance.AcceptedBy.Organization = eulaOrg.Text()
			FA.EulaAcceptance.AcceptedBy.JobTitle = eulaTitle.Text()

			//marshal (json encode) the map into a json string
			FAData, err := json.Marshal(FA)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//make the rest call with the json payload and stores response
			resp := apiCallFB("PATCH", "http://"+tempIP.Text()+":8081/array-initial-config", "", FAData)
			//testing
			//resp := apiCallFB("PATCH", "https://pureapisim.azurewebsites.net/api/array-initial-config", "", FAData)
			//update the initResult field with response.
			if statusCode == 200 {
				initResult.SetText("Success! \n\nResponse: \n\n" + string(resp))
			} else {
				initResult.SetText("Error! \n\nStatus Code: \n" + string(statusCode) + "\n\nResponse:\n" + string(resp))
			}

		}
	})

	return hbox
}

//MAIN FLASHBLADE UI WORKER FUNCTION//
func initializeFBTab() ui.Control {
	//results field variable used throughout as a "console out"
	initResult := ui.NewMultilineEntry()

	//first column definition
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	//define vertical box inside column similar to a div
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	//BUTTONS GROUP FOR LEFT COLUMN - FORM CONTROLS//
	//define the group for the form
	buttonGroup := ui.NewGroup("Form Controls")
	buttonGroup.SetMargined(true)

	//add group to the vertical box
	vbox.Append(buttonGroup, false)

	///Form Instantiation///
	//define the form for the button group
	buttonForm := ui.NewForm()
	buttonForm.SetPadded(true)

	///Button Definition Login///
	//embed the login form field inside the first form group
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("STEP 1 Login", ui.NewLabel(""), false)
	login := ui.NewButton("Login Page")
	buttonForm.Append("Login Form", login, false)
	//seperator line
	hbox.Append(ui.NewVerticalSeparator(), false)
	///End Button Definition///

	///Button Definition Array///
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("STEP 2 Array Config", ui.NewLabel(""), false)
	array := ui.NewButton("Array Form")
	buttonForm.Append("Array Form", array, false)
	///End Button Definition///

	///Button Definition DNS///
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("STEP 3 DNS Config", ui.NewLabel(""), false)
	dns := ui.NewButton("DNS Form")
	buttonForm.Append("DNS Form", dns, false)
	///End Button Definition///

	///Button Definition Subnets Aggregation///
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("STEP 4 Subnet Config", ui.NewLabel(""), false)
	subnet := ui.NewButton("Subnet Form")
	buttonForm.Append("Subnet Form", subnet, false)
	///End Button Definition///

	///Button Definition Network Interfaces///
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("STEP 5 Network Config", ui.NewLabel(""), false)
	network := ui.NewButton("NIC Form")
	buttonForm.Append("NIC Form", network, false)
	///End Button Definition///

	///Button Definition smtp///
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("STEP 6 SMTP Config", ui.NewLabel(""), false)
	smtp := ui.NewButton("SMTP Form")
	buttonForm.Append("SMTP Form", smtp, false)
	///End Button Definition///

	///Button Definition support///
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("STEP 7 Support Config", ui.NewLabel(""), false)
	support := ui.NewButton("Support Form")
	buttonForm.Append("Phonehome Form", support, false)
	///End Button Definition///

	///Button Definition alert watchers///
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("STEP 8 Alerts Config", ui.NewLabel(""), false)
	aw := ui.NewButton("Alerts Form")
	buttonForm.Append("Alerts Form", aw, false)
	///End Button Definition///

	///Button Definition validation and finalization///
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("STEP 9 Final Step", ui.NewLabel(""), false)
	final := ui.NewButton("Finalize Form")
	buttonForm.Append("Finalize Form", final, false)
	///End Button Definition///

	///Button Definition validation and advanced///
	buttonGroup.SetChild(buttonForm)
	buttonForm.Append("", ui.NewLabel(""), false)
	advanced := ui.NewButton("Advanced")
	buttonForm.Append("Advanced Options", advanced, false)
	///End Button Definition///

	//Middle column
	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	//Login FORM//
	loginGroup := ui.NewGroup("Login")
	loginGroup.SetMargined(false)
	vbox.Append(loginGroup, false)
	loginForm := ui.NewForm()
	loginForm.SetPadded(true)
	loginGroup.SetChild(loginForm)
	loginGroup.Hide()
	//variables
	//TESTING ONLY//
	// apiUrlForm.SetText("https://pureapisim.azurewebsites.net/api/1.8.1")
	// apiToken.SetText("2PDoD5iaokKDwGh9uNqt1jpDTNpgshfiOzO643z5ch92Mwycl7veBA==")
	//END TESTING ONLY//
	apiToken := ui.NewEntry()
	apiToken.SetText("PURESETUP")
	xAuthTokenField := ui.NewEntry()
	loginSubmitButton := ui.NewButton("Create Session")
	getAPIVersionsButton := ui.NewButton("Generate URL")
	apiUrlForm := ui.NewEntry()
	managementIP := ui.NewEntry()
	//append variables to form
	loginForm.Append("Array API URL", apiUrlForm, false)
	loginForm.Append("", ui.NewLabel("format:  https://10.1.1.100/api/1.8"), false)
	loginForm.Append("", loginSubmitButton, false)
	loginForm.Append("", ui.NewLabel(""), false)
	loginForm.Append("", ui.NewLabel("     __________OR__________"), false)
	loginForm.Append("", ui.NewLabel("  AUTO GENERATE THE API URL"), false)
	loginForm.Append("IP of FB", managementIP, false)
	loginForm.Append("Generate URL", getAPIVersionsButton, false)

	//Array Form//
	//variables
	arrayName := ui.NewEntry()
	ntpServer := ui.NewEntry()
	timeZone := ui.NewCombobox()
	tz := timeZones()
	for i, v := range tz {
		timeZone.Append(v)
		i++
	}
	//define the form
	arrayGroup := ui.NewGroup("Array Config")
	arrayGroup.SetMargined(false)
	vbox.Append(arrayGroup, false)
	arrayGroup.Hide()
	arrayForm := ui.NewForm()
	arrayForm.SetPadded(true)
	arrayGroup.SetChild(arrayForm)
	arrayForm.Append("Array Name", arrayName, false)
	arrayForm.Append("NTP Servers", ntpServer, false)
	arrayForm.Append("", ui.NewLabel("*Comma seperated for multiple entries"), false)
	arrayForm.Append("TimeZone", timeZone, false)
	arrayGetButton := ui.NewButton("Query Array")
	arrayPatchButton := ui.NewButton("Apply To Array")
	arrayForm.Append("", arrayPatchButton, false)
	arrayForm.Append("", ui.NewLabel(""), false)
	arrayForm.Append("", arrayGetButton, false)
	//end Array Form//

	//DNS Form//
	//variables
	dnsDomain := ui.NewEntry()
	dnsServer := ui.NewEntry()
	//define the form
	dnsGroup := ui.NewGroup("DNS Config")
	dnsGroup.SetMargined(false)
	vbox.Append(dnsGroup, false)
	dnsGroup.Hide()
	dnsForm := ui.NewForm()
	dnsForm.SetPadded(true)
	dnsGroup.SetChild(dnsForm)
	dnsForm.Append("DNS Domain Name", dnsDomain, false)
	dnsForm.Append("DNS Servers", dnsServer, false)
	dnsForm.Append("", ui.NewLabel("*Comma seperated for multiple entries"), false)
	dnsGetButton := ui.NewButton("Query Array")
	dnsPatchButton := ui.NewButton("Apply To Array")
	dnsForm.Append("", dnsPatchButton, false)
	dnsForm.Append("", ui.NewLabel(""), false)
	dnsForm.Append("", dnsGetButton, false)
	//end DNS Form//

	//SHOWN IN ADVANCED SECTION//
	//LAG display Buttons to show sub-forms/
	lagNew := ui.NewButton("Create New LAG")
	lagExisting := ui.NewButton("Update Existing")
	lagGetButton := ui.NewButton("Query LAG")
	lagDelete := ui.NewButton("Delete LAG")
	lagGroupInit := ui.NewGroup("LAG Options")
	lagGroupInit.SetMargined(false)
	vbox.Append(lagGroupInit, false)
	lagGroupInit.Hide()
	lagFormInit := ui.NewForm()
	lagFormInit.SetPadded(true)
	lagGroupInit.SetChild(lagFormInit)
	lagFormInit.Append("", lagNew, false)
	lagFormInit.Append("", lagExisting, false)
	lagFormInit.Append("", lagDelete, false)
	lagFormInit.Append("", lagGetButton, false)

	//lag create new group and form
	lagNameNew := ui.NewEntry()
	lagNameExisting := ui.NewEntry()
	lagPortsNew := ui.NewEntry()
	lagPortsExisting := ui.NewEntry()
	lagAddRemove := ui.NewCombobox()
	lagAddRemove.Append("Add Ports")
	lagAddRemove.Append("Remove Ports")
	lagGroupNew := ui.NewGroup("New LAG Config")
	lagGroupNew.SetMargined(false)
	vbox.Append(lagGroupNew, false)
	lagGroupNew.Hide()
	lagFormNew := ui.NewForm()
	lagFormNew.SetPadded(true)
	lagGroupNew.SetChild(lagFormNew)
	lagFormNew.Append("LAG Name", lagNameNew, false)
	lagFormNew.Append("Lag Port Name(s)", lagPortsNew, false)
	lagFormNew.Append("", ui.NewLabel("E.g. CH1.FM1.ETH1..."), false)
	lagPostButton := ui.NewButton("Create New LAG")
	lagFormNew.Append("", lagPostButton, false)

	//lag modify existing group and form
	lagGroupExisting := ui.NewGroup("Existing LAG Config")
	lagGroupExisting.SetMargined(false)
	vbox.Append(lagGroupExisting, false)
	lagGroupExisting.Hide()
	lagFormExisting := ui.NewForm()
	lagFormExisting.SetPadded(true)
	lagGroupExisting.SetChild(lagFormExisting)

	lagFormNew.Append("", ui.NewLabel(""), false)
	lagFormExisting.Append("", ui.NewLabel(""), false)
	lagFormExisting.Append("LAG Name", lagNameExisting, false)
	lagFormExisting.Append("Lag Port Name(s)", lagPortsExisting, false)
	lagFormExisting.Append("", ui.NewLabel("*Comma seperated for multiple entries"), false)

	lagFormExisting.Append("", ui.NewLabel(""), false)
	lagFormExisting.Append("Modify Ports", lagAddRemove, false)
	lagPatchButton := ui.NewButton("Update LAG Ports")
	lagFormExisting.Append("", lagPatchButton, false)

	//lag create delete group and form
	lagNameDelete := ui.NewEntry()
	lagDeleteConfirm := ui.NewCheckbox("Yes")
	lagGroupDelete := ui.NewGroup("LAG Delete")
	lagGroupDelete.SetMargined(false)
	vbox.Append(lagGroupDelete, false)
	lagGroupDelete.Hide()
	lagFormDelete := ui.NewForm()
	lagFormDelete.SetPadded(true)
	lagGroupDelete.SetChild(lagFormDelete)
	lagFormDelete.Append("LAG Name", lagNameDelete, false)
	lagFormDelete.Append("Confirm Delete", lagDeleteConfirm, false)
	lagDeleteButton := ui.NewButton("Delete LAG")
	lagFormDelete.Append("", lagDeleteButton, false)
	//END link aggrigation Form//
	//END ADVANCED SECTION//

	//subnets Form//
	subnetGateway := ui.NewEntry()
	subnetLag := ui.NewEntry()
	subnetLag.SetText("")
	subnetMtu := ui.NewEntry()
	subnetMtu.SetText("1500")
	subnetPrefix := ui.NewEntry()
	subnetVlan := ui.NewEntry()
	subnetVlan.SetText("0")
	subnetName := ui.NewEntry()
	subnetName.SetText("mgmt")
	subnetOOB := ui.NewCombobox()
	subnetOOB.Append("true")
	subnetOOB.Append("false")
	subnetGroup := ui.NewGroup("Subnet Config")
	subnetGroup.SetMargined(false)
	vbox.Append(subnetGroup, false)
	subnetGroup.Hide()
	subnetForm := ui.NewForm()
	subnetForm.SetPadded(true)
	subnetGroup.SetChild(subnetForm)
	subnetForm.Append("Subnet Name", subnetName, false)
	subnetForm.Append("Gateway IP", subnetGateway, false)
	subnetForm.Append("Subnet Prefix", subnetPrefix, false)
	subnetForm.Append("", ui.NewLabel("Prefix e.g. 10.1.1.0/24"), false)
	subnetForm.Append("VLAN", subnetVlan, false)
	subnetForm.Append("Out of Band", subnetOOB, false)
	subnetGetButton := ui.NewButton("Query")
	subnetPatchButton := ui.NewButton("Update Existing")
	subnetPostButton := ui.NewButton("Create New")
	subnetDeleteButton := ui.NewButton("Delete")
	subnetForm.Append("", subnetPostButton, false)
	subnetForm.Append("", subnetPatchButton, false)
	subnetForm.Append("", subnetGetButton, false)
	subnetForm.Append("", subnetDeleteButton, false)
	//end subnets Form//

	//network interfaces Form//
	virIP := ui.NewEntry()
	fm1AdminIP := ui.NewEntry()
	fm2AdminIP := ui.NewEntry()
	nicGroup := ui.NewGroup("Net Interface Config")
	nicGroup.SetMargined(false)
	vbox.Append(nicGroup, false)
	nicGroup.Hide()
	nicForm := ui.NewForm()
	nicForm.SetPadded(true)
	nicGroup.SetChild(nicForm)

	nicForm.Append("", ui.NewLabel("Admin VIR0"), false)
	nicForm.Append("IP Address", virIP, false)

	nicForm.Append("", ui.NewLabel(""), false)
	nicForm.Append("", ui.NewLabel("Admin FM1"), false)
	nicForm.Append("IP Address", fm1AdminIP, false)

	nicForm.Append("", ui.NewLabel(""), false)
	nicForm.Append("", ui.NewLabel("Admin FM2"), false)
	nicForm.Append("IP Address", fm2AdminIP, false)

	nicGetButton := ui.NewButton("Query Array")
	nicPatchButton := ui.NewButton("Apply NIC Config")

	nicForm.Append("", nicPatchButton, false)
	nicForm.Append("", nicGetButton, false)
	//end network interfaces Form//

	//smtp Form//
	smtpRelayHost := ui.NewEntry()
	smtpSenderDomain := ui.NewEntry()

	smtpGroup := ui.NewGroup("SMTP Config")
	smtpGroup.SetMargined(false)
	vbox.Append(smtpGroup, false)
	smtpGroup.Hide()
	smtpForm := ui.NewForm()
	smtpForm.SetPadded(true)
	smtpGroup.SetChild(smtpForm)
	smtpForm.Append("Sender Domain", smtpSenderDomain, false)
	smtpForm.Append("Relay Host (optional)", smtpRelayHost, false)
	smtpGetButton := ui.NewButton("Query")
	smtpPatchButton := ui.NewButton("Create New")
	smtpForm.Append("", smtpPatchButton, false)
	smtpForm.Append("", smtpGetButton, false)
	//end smtp Form//

	//support Form//
	supportPhoneHome := ui.NewCombobox()
	supportPhoneHome.Append("true")
	supportPhoneHome.Append("false")
	supportProxy := ui.NewEntry()
	supportGroup := ui.NewGroup("Support Config")
	supportGroup.SetMargined(false)
	vbox.Append(supportGroup, false)
	supportGroup.Hide()
	supportForm := ui.NewForm()
	supportForm.SetPadded(true)
	supportGroup.SetChild(supportForm)
	supportForm.Append("Enable Phone Home?", supportPhoneHome, false)
	supportForm.Append("Proxy Server (optional)", supportProxy, false)
	supportGetButton := ui.NewButton("Query Array")
	supportPatchButton := ui.NewButton("Apply To Array")
	supportForm.Append("", supportPatchButton, false)
	supportForm.Append("", supportGetButton, false)
	//end support Form//

	//alert watchers Form//
	awName := ui.NewEntry()
	awEnabled := ui.NewCombobox()
	awEnabled.Append("true")
	awEnabled.Append("false")
	awGroup := ui.NewGroup("Alert Watchers Config")
	awGroup.SetMargined(false)
	vbox.Append(awGroup, false)
	awGroup.Hide()
	awForm := ui.NewForm()
	awForm.SetPadded(true)
	awGroup.SetChild(awForm)
	awForm.Append("Email Address", awName, false)
	awForm.Append("Enabled", awEnabled, false)
	awGetButton := ui.NewButton("Query")
	awPatchButton := ui.NewButton("Update Existing")
	awDeleteButton := ui.NewButton("Delete Alert Watcher")
	awPostButton := ui.NewButton("New Alert Watcher")
	awForm.Append("", awPostButton, false)
	awForm.Append("", awPatchButton, false)
	awForm.Append("", awDeleteButton, false)
	awForm.Append("", awGetButton, false)
	//end alert watchers Form//

	//finalization and validation Form//
	finalSetupComplete := ui.NewCombobox()
	finalSetupComplete.Append("true")
	finalSetupComplete.Append("false")
	finalGroup := ui.NewGroup("Validate and Finalize")
	finalGroup.SetMargined(false)
	vbox.Append(finalGroup, false)
	finalGroup.Hide()
	finalForm := ui.NewForm()
	finalForm.SetPadded(true)
	finalGroup.SetChild(finalForm)
	finalForm.Append("Setup Complete", finalSetupComplete, false)
	finalGetButton := ui.NewButton("Validate")
	finalPatchButton := ui.NewButton("Finalize Setup")
	finalForm.Append("", finalPatchButton, false)
	finalForm.Append("", finalGetButton, false)
	//end finalization and validation Form//

	//vertical seperator between 2nd and 3rd column.
	hbox.Append(ui.NewVerticalSeparator(), false)

	//THIRD COLUMN DEFINITION//
	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, true)

	//Innformation group for output//
	group9 := ui.NewGroup("Array Information")
	group9.SetMargined(true)
	vbox.Append(group9, true)

	entryForm9 := ui.NewForm()
	entryForm9.SetPadded(true)
	group9.SetChild(entryForm9)

	//labels used to display the api url and x-auth token in third column
	xAuthTokenLabel := ui.NewLabel("")
	apiUrlLabel := ui.NewLabel("")
	xAuthTokenField.SetReadOnly(true)
	prog := ui.NewProgressBar()
	prog.SetValue(0)
	entryForm9.Append("API URL: ", apiUrlLabel, false)
	entryForm9.Append("X-Auth-Token", xAuthTokenLabel, false)
	//progress bar
	entryForm9.Append("Progress: ", prog, false)

	//multiline field for showing results of patch api call and form validation messages.
	//sets the initResults console to readonly
	initResult.SetReadOnly(true)
	entryForm9.Append("Init Results", initResult, true)
	initResult.SetText("Welcome to the FlashBlade Zero Touch Provisioner.\n\nThe tool has 3 colunms.  Form Controls, Form Actions, and Information.\nStart with Step 1 on the left and proceed through Step 9 to finish the Array config.\n\nThe progress bar above ^^^ will fill as you proceed through the steps.\n\nFor questions or assistance please reach out to your Pure account team.\n\nThank you! ")

	//Login Form Button
	login.OnClicked(func(*ui.Button) {
		loginGroup.Show()
		arrayGroup.Hide()
		dnsGroup.Hide()

		lagGroupNew.Hide()
		lagGroupExisting.Hide()
		lagGroupInit.Hide()
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("Start here and  fill out this form to logon the the Array to proceed.\n\nMore Info:\nThe Array API URL should be in the format of:\nhttps://<FB DHCP IP>/api/<api version>\n\nGenerate URL section:\nYou can enter the DHCP IP of the FB array into the Auto Generate section and the tool will build the API URL in the correct format for you.  You will still need to provide the API Token.")

	})

	//arrays Form Button
	array.OnClicked(func(*ui.Button) {
		arrayGroup.Show()
		loginGroup.Hide()
		dnsGroup.Hide()

		lagGroupNew.Hide()
		lagGroupExisting.Hide()
		lagGroupInit.Hide()
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("Step 2.  Provide the name for this Array.  Some rules on that are as follows: The Array Name cannot begin or end with a dash (but CAN contain a dash).  2. the name cannot exceed 55 characters in length.\n\nEnter your NTP Server or Servers.  If you have more than one to enter, plase seperate them by commas with no spaces.\n\nYou can also Query the array before and after to see the status of this section.")

	})

	//DNS Form Button
	dns.OnClicked(func(*ui.Button) {
		dnsGroup.Show()
		arrayGroup.Hide()
		loginGroup.Hide()

		lagGroupNew.Hide()
		lagGroupExisting.Hide()
		lagGroupInit.Hide()
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("Step 3.  Provide the DNS Domain Name for your environment as well as at least 1 (2 recommended) DNS server.  If you have more than one to enter, plase seperate them by commas with no spaces.\n\nYou can also Query the array before and after to see the status of this section.")

	})

	//LAG New Form Button
	lagNew.OnClicked(func(*ui.Button) {
		lagGroupInit.Show()
		lagGroupNew.Show()
		lagGroupExisting.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupDelete.Hide()

		initResult.SetText("Advanced Section, Create new LAG.  Provide the new LAG Name and Port Names. If you have more than one to enter, plase seperate them by commas with no spaces.\n\nYou can also Query the array before and after to see the status of this section.")

	})

	//LAG existing Form Button
	lagExisting.OnClicked(func(*ui.Button) {
		lagGroupInit.Show()
		lagGroupExisting.Show()
		lagGroupNew.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupDelete.Hide()

		initResult.SetText("Advanced Section, Update LAG Ports.  You must provide a valid LAG Name the enter the Port Names you wish to change.  Finally select if you want to add or remove these Ports.  If you have more than one to enter, plase seperate them by commas with no spaces.\n\nYou can also Query the array before and after to see the status of this section.")

	})

	//LAG delete Form Button
	lagDelete.OnClicked(func(*ui.Button) {
		lagGroupDelete.Show()
		lagGroupInit.Show()
		lagGroupExisting.Hide()
		lagGroupNew.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()

		initResult.SetText("Advanced Section, Delete LAG.  Provide the exisiting LAG Name to delete.")

	})

	//Subnet Form Button
	subnet.OnClicked(func(*ui.Button) {
		finalGroup.Hide()
		subnetGroup.Show()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupInit.Hide()
		lagGroupExisting.Hide()
		lagGroupNew.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("Step 4. Please provide the Gateway IP for this Array as well as the Subnet Prefix in teh format of x.x.x.x/x. e.g. 10.1.1.0/24.\n\nThe Subnet Name and VLAN are pre-populated for you as the most common names.  You can change but make sure you have a reason for doing so.\n\nYou can also Query the array before and after to see the status of this section.")

	})

	//Network Init Form Button
	network.OnClicked(func(*ui.Button) {
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Show()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupInit.Hide()
		lagGroupExisting.Hide()
		lagGroupNew.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("Step 5. Provide the individual IP address for the Virtual 0 Nic, FM1 NIC and FM2 NIC.\n\nNote, this step takes 10 -20 seconds to complete.\n\nYou can also Query the array before and after to see the status of this section.")

	})

	//SMTP Form Button
	smtp.OnClicked(func(*ui.Button) {
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Show()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupInit.Hide()
		lagGroupExisting.Hide()
		lagGroupNew.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("Step 6.  Provide the email domain for your environment.  e.g. example.com.\n\nOptionally, you may provide the SMTP relay host to use for outbound email from the Array.\n\nYou can also Query the array before and after to see the status of this section.")

	})

	//Support Form Button
	support.OnClicked(func(*ui.Button) {
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Show()
		awGroup.Hide()
		lagGroupInit.Hide()
		lagGroupExisting.Hide()
		lagGroupNew.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("Step 7.  Select whether or not to enable Phone Home for integration with Pure1.  If you are enabling and have a proxy server for outbound connectivity, please enter it as well.\n\nYou can also Query the array before and after to see the status of this section.")

	})

	//Alert Watchers Form Button
	aw.OnClicked(func(*ui.Button) {
		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Show()
		lagGroupInit.Hide()
		lagGroupExisting.Hide()
		lagGroupNew.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("Step 8. Provide a valid email address will recieve all alerts generated by the Array.  Please enter one email addresse at a time.\n\nIf needed you can modify or delete email address.\n\nYou can also Query the array before and after to see the status of this section.")

	})

	//Final Validation Form Button
	final.OnClicked(func(*ui.Button) {
		finalGroup.Show()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupInit.Hide()
		lagGroupExisting.Hide()
		lagGroupNew.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("Final Step.  When all previous steps are complete you can Validate and then Finalize the Config.\n\nValidate should return 'true' in the response.\n\nOnce Finalized you will no longer be able to access the array using this tool and just simply exit using the X at the top right.")

	})

	//Advanced Form Button
	advanced.OnClicked(func(*ui.Button) {
		lagGroupInit.Show()

		finalGroup.Hide()
		subnetGroup.Hide()
		nicGroup.Hide()
		smtpGroup.Hide()
		supportGroup.Hide()
		awGroup.Hide()
		lagGroupExisting.Hide()
		lagGroupNew.Hide()

		dnsGroup.Hide()
		arrayGroup.Hide()
		loginGroup.Hide()
		lagGroupDelete.Hide()
		initResult.SetText("LAG Config. Use only with assistnace of Pure Storage Support.")

	})

	//Buttons Actions from Forms//
	//QUERY FOR API VERSIONS//
	getAPIVersionsButton.OnClicked(func(*ui.Button) {
		//make sure the api endpoints are in the right format
		passed := true
		//Validate inputs
		validate := validator.New()
		err := validate.Var(managementIP.Text(), "required,ipv4")
		if err != nil {
			initResult.SetText("Please provide a valid IP Address for the FB management endpoint")
			passed = false

		}
		//if passed validation
		if passed == true {
			//make the rest call
			resp := apiCallFB("GET", "https://"+managementIP.Text()+"/api/api_version", apiToken.Text(), nil)

			type Version struct {
				Versions []string `json:"versions"`
			}

			var version Version
			err := json.Unmarshal(resp, &version)
			if err == nil {
				fmt.Println(err)
			}
			fmt.Printf("%v", (len(version.Versions)))
			if len(version.Versions) > 0 {
				apiUrlForm.SetText("https://" + managementIP.Text() + "/api/" + version.Versions[(len(version.Versions)-1)])
				apiUrlLabel.SetText("https://" + managementIP.Text() + "/api/" + version.Versions[(len(version.Versions)-1)])
			}
			//set the response in the display of the app
			initResult.SetText(string(resp))

		}

	})

	//LOGIN SUBMIT//
	loginSubmitButton.OnClicked(func(*ui.Button) {
		//returns a slice broken out by forward slash in the url
		url := strings.Split(apiUrlForm.Text(), "/")
		//make sure the api endpoints are in the right format
		passed := true
		if len(url) > 4 {
			loginUrl = (url[0] + "//" + url[2] + "/" + url[3])
			apiUrl = (url[0] + "//" + url[2] + "/" + url[3] + "/" + url[4])
		} else {
			initResult.SetText("please enter a valid API url that includes the version.  e.g. https://purefb01.example.com/api/1.8")
			passed = false
		}
		if passed == true {
			apiUrlLabel.SetText(apiUrl)
			fmt.Println(apiUrl)

			//make the rest call
			resp := postAPICallLoginFB(loginUrl+"/login", apiToken.Text())

			xAuthTokenField.SetText(xAuthToken)
			xAuthTokenLabel.SetText(xAuthToken)

			//check if the post was a success
			if statusCode == 200 {
				//used for progress bar
				if t, found := steps[1]; found {
					fmt.Println("step ", t, " already completed not changing progress bar.")
				} else {
					steps[1] = "complete"
					progressCounter = progressCounter + 10
					prog.SetValue(progressCounter)
				}
				//set the response in the display of the app
				initResult.SetText(string(resp) + "\n\nLogon Successful!\n\nPlease proceed to the Array form.")
			}
		}
	})

	//action for array Get button to make api call
	arrayGetButton.OnClicked(func(*ui.Button) {
		result := apiCallFB("GET", apiUrl+"/arrays", xAuthToken, nil)
		initResult.SetText(string(result))
	})

	//action for array Apply button to make api call
	arrayPatchButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate Array Name
		var rxPat = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,54}[a-zA-Z0-9])?$`)
		if !rxPat.MatchString(arrayName.Text()) {
			initResult.SetText("Array Name has blank or contains invalid characters.  It must begin with a number or letter, can contain a dash in the body of the name, but must also end with a number or letter.   No more than 55 characters in length.")
			passed = false
		}
		//validate ntp server or servers
		err1 := validate.Var(ntpServer.Text(), "required")
		if err1 != nil {
			initResult.SetText("Please provide the NTP server(s)")
			passed = false
		}
		ntp := strings.Split(ntpServer.Text(), ",")
		for i := 0; i < len(ntp); i++ {
			//fmt.Print(ntp[i] + "\n")
			err2 := validate.Var(ntp[i], "fqdn|ip")
			if err2 != nil {
				initResult.SetText("Please provide a fqdn or ip for the NTP server.\n\nIf more than one server is entered please use comma seperation with no spaces in-between.")
				passed = false
			}
		}
		if timeZone.Selected() < 0 {
			initResult.SetText("Please select a Timezone")
			passed = false
		}
		if passed == true {
			//struct here
			type FAB struct {
				Name      string   `json:"name"`
				NtpServer []string `json:"ntp_servers"`
				TimeZone  string   `json:"time_zone"`
			}

			//initialize FAS struct object
			FB := &FAB{}
			FB.Name = arrayName.Text()
			FB.NtpServer = ntp
			FB.TimeZone = tz[timeZone.Selected()]

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			//send patch request to array
			result := apiCallFB("PATCH", apiUrl+"/arrays", xAuthToken, FBData)

			//check if request was successful (code 200)
			if statusCode == 200 {
				//decode the response from a successfull call to a map interface//
				var r map[string]interface{}
				json.Unmarshal([]byte(result), &r)
				name, d, o, err := jsonparser.Get(result, "name")
				fmt.Print("jsonparser out: ", d, o, err)
				ntpservers, d, o, err := jsonparser.Get(result, "ntp_servers")
				fmt.Print("jsonparser out: ", d, o, err)
				timezone, d, o, err := jsonparser.Get(result, "time_zone")
				fmt.Print("jsonparser out: ", d, o, err)
				fmt.Print(string(name))
				initResult.SetText("Success!\n\nApplied the following:\n\nName: " + string(name) + "\nNTP Servers: " + string(ntpservers) + "\nTimeZone: " + string(timezone) + "\n\nPlease proceed to the DNS form.")

				//used for the progress bar
				//check if post was success and if this step has previously been completed.
				if t, found := steps[2]; found {
					fmt.Println("step ", t, ", already completed not changing progress bar.")
				} else {
					steps[2] = "complete"
					progressCounter = progressCounter + 10
					prog.SetValue(progressCounter)
				}
			} else {
				initResult.SetText(string(result))
			}

		}
	})

	dnsGetButton.OnClicked(func(*ui.Button) {
		result := apiCallFB("GET", apiUrl+"/dns", xAuthToken, nil)
		initResult.SetText(string(result))
	})

	dnsPatchButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate DNS server entries
		err := validate.Var(dnsServer.Text(), "required")
		if err != nil {
			initResult.SetText("Please provide the DNS Server(s)")
			passed = false
		}
		//split multiple entries into a string array
		dns := strings.Split(dnsServer.Text(), ",")
		for i := 0; i < len(dns); i++ {
			err1 := validate.Var(dns[i], "fqdn|ip")
			if err1 != nil {
				initResult.SetText("Please provide a fqdn or ip for the DNS server.\n\nIf more than one server is entered please use comma seperation with no spaces in-between.")
				passed = false
			}
		}
		//validate dns domain
		err2 := validate.Var(dnsDomain.Text(), "fqdn")
		if err2 != nil {
			initResult.SetText("Please provide the DNS Domain")
			passed = false
		}
		if passed == true {

			//struct here
			type FAB struct {
				Domain      string   `json:"domain"`
				Nameservers []string `json:"nameservers"`
			}

			//initialize FAS struct object
			FB := &FAB{}
			FB.Domain = dnsDomain.Text()
			FB.Nameservers = dns

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			//make the Patch rest api call to the array
			result := apiCallFB("PATCH", apiUrl+"/dns", xAuthToken, FBData)
			//check if success to provide results
			if statusCode == 200 {
				//decode the response from a successfull call to a map interface//
				var r map[string]interface{}
				json.Unmarshal([]byte(result), &r)
				domain, d, o, err := jsonparser.Get(result, "domain")
				fmt.Print("jsonparser out: ", d, o, err)
				nameservers, d, o, err := jsonparser.Get(result, "nameservers")
				fmt.Print("jsonparser out: ", d, o, err)
				//print results to output field
				initResult.SetText("Success!\n\nApplied the following:\n\nDomain: " + string(domain) + "\nName Servers: " + string(nameservers) + "\n\nPlease proceed to the Subnets Form.")

				//used for the progress bar
				//check if post was success and if this step has previously been completed.
				if t, found := steps[3]; found {
					fmt.Println("step ", t, ", already completed not changing progress bar.")
				} else {
					steps[2] = "complete"
					progressCounter = progressCounter + 10
					prog.SetValue(progressCounter)
				}

			} else {
				initResult.SetText("name: " + string(result))
			}
		}
	})

	//Lag Buttons
	lagGetButton.OnClicked(func(*ui.Button) {

		result := apiCallFB("GET", apiUrl+"/link-aggregation-groups", xAuthToken, nil)
		initResult.SetText(string(result))

	})

	lagPostButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate lag port and  server name
		err2 := validate.Var(lagNameNew.Text(), "required")
		if err2 != nil {
			initResult.SetText("Please provide the Array name")
			passed = false
		}
		err1 := validate.Var(lagPortsNew.Text(), "required")
		if err1 != nil {
			initResult.SetText("Please provide the Port Name(s)")
			passed = false
		}

		//if validation passed
		if passed == true {

			//manually build the post request
			portNames := strings.Split(lagPortsNew.Text(), ",")
			var pName = `{"ports": [`
			for i, v := range portNames {
				i++
				pName += `{"name": "`
				pName += v
				pName += `"}`
				if i < len(portNames) {
					pName += `,`
				}
			}
			pName += `]}`
			pNameSlice := []byte(pName)

			result := apiCallFB("POST", apiUrl+"/link-aggregation-groups?names="+lagNameNew.Text(), xAuthToken, pNameSlice)
			initResult.SetText(string(result))
		}
	})

	lagPatchButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate
		err2 := validate.Var(lagNameExisting.Text(), "required")
		if err2 != nil {
			initResult.SetText("Please provide the Array name")
			passed = false
		}

		if passed == true {

			//manually build the portNames JSON structure.
			portNames := strings.Split(lagPortsExisting.Text(), ",")
			var pName = ""
			if lagAddRemove.Selected() == 0 {
				pName += `{"add_ports":[`
			}
			if lagAddRemove.Selected() == 1 {
				pName += `{"remove_ports":[`
			}
			for i, v := range portNames {
				i++
				pName += `{"name":"`
				pName += v
				pName += `"}`
				if i < len(portNames) {
					pName += `,`
				}
			}
			pName += `]}`
			pNameSlice := []byte(pName)

			result := apiCallFB("PATCH", apiUrl+"/link-aggregation-groups?names="+lagNameExisting.Text(), xAuthToken, pNameSlice)
			initResult.SetText(string(result))
		}
	})

	lagDeleteButton.OnClicked(func(*ui.Button) {

		result := apiCallFB("DELETE", apiUrl+"/link-aggregation-groups?names="+lagNameDelete.Text(), xAuthToken, nil)
		initResult.SetText(string(result))

	})

	subnetGetButton.OnClicked(func(*ui.Button) {
		result := apiCallFB("GET", apiUrl+"/subnets", xAuthToken, nil)
		initResult.SetText(string(result))
	})

	subnetPostButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate
		err1 := validate.Var(subnetVlan.Text(), "numeric,min=1")
		if err1 != nil {
			initResult.SetText("Please provide a valid vlanID")
			passed = false
		}
		err2 := validate.Var(subnetGateway.Text(), "ip")
		if err2 != nil {
			initResult.SetText("Please provide a valid gateway IP")
			passed = false
		}
		err4 := validate.Var(subnetPrefix.Text(), "cidr")
		if err4 != nil {
			initResult.SetText("Please provide the Prefix")
			passed = false
		}
		err5 := validate.Var(subnetName.Text(), "required")
		if err5 != nil {
			initResult.SetText("Please provide the Subnet Name")
			passed = false
		}
		if subnetOOB.Selected() < 0 {
			initResult.SetText("Please select true or false for Out of Band")
			passed = false
		}

		if passed == true {

			type FAB struct {
				Gateway string `json:"gateway"`
				Mtu     int    `json:"mtu"`
				Prefix  string `json:"prefix"`
				Vlan    int    `json:"vlan"`
				OOB     bool   `json:"out_of_band"`
			}

			//set out of band value for post request
			var OOB bool
			if subnetOOB.Selected() == 0 {
				OOB = true
			}
			if subnetOOB.Selected() == 1 {
				OOB = false
			}

			//convert the vlan entry from string to int.
			vlanInt, err := strconv.Atoi(subnetVlan.Text())

			//initialize FAS struct object
			FB := &FAB{}
			FB.OOB = OOB
			FB.Gateway = subnetGateway.Text()
			FB.Mtu = 1500
			FB.Prefix = subnetPrefix.Text()
			FB.Vlan = vlanInt

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			result := apiCallFB("POST", apiUrl+"/subnets?names="+subnetName.Text(), xAuthToken, FBData)
			if statusCode == 200 {
				//decode the response from a successfull call to a map interface//
				var r map[string]interface{}
				json.Unmarshal([]byte(result), &r)
				name, d, o, err := jsonparser.Get(result, "items", "[0]", "name")
				fmt.Print("jsonparser out: ", d, o, err)
				gateway, d, o, err := jsonparser.Get(result, "items", "[0]", "gateway")
				prefix, d, o, err := jsonparser.Get(result, "items", "[0]", "prefix")
				vlan, d, o, err := jsonparser.Get(result, "items", "[0]", "vlan")
				fmt.Print("jsonparser out: ", d, o, err)

				//print results to output field
				initResult.SetText("Success!\n\nApplied the following:\n\nSubnetName: " + string(name) + "\nGateway: " + string(gateway) + "\nPrefix: " + string(prefix) + "\nVLAN: " + string(vlan) + "\n\nPlease Proceed to the NIC Form")

				//used for the progress bar to check if this step has previously been completed.
				if t, found := steps[4]; found {
					fmt.Println("step ", t, ", already completed not changing progress bar.")
				} else {
					steps[4] = "complete"
					progressCounter = progressCounter + 10
					prog.SetValue(progressCounter)
				}
			} else {
				initResult.SetText(string(result))
			}
		}
	})

	subnetPatchButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate entries
		err1 := validate.Var(subnetVlan.Text(), "numeric,min=1")
		if err1 != nil {
			initResult.SetText("Please provide a valid vlanID")
			passed = false
		}
		err2 := validate.Var(subnetGateway.Text(), "ip")
		if err2 != nil {
			initResult.SetText("Please provide a valid gateway IP")
			passed = false
		}
		err4 := validate.Var(subnetPrefix.Text(), "required")
		if err4 != nil {
			initResult.SetText("Please provide the Prefix")
			passed = false
		}
		err5 := validate.Var(subnetName.Text(), "required")
		if err5 != nil {
			initResult.SetText("Please provide the Subnet Name")
			passed = false
		}
		if subnetOOB.Selected() < 0 {
			initResult.SetText("Please select true or false for Out of Band")
			passed = false
		}

		if passed == true {

			//defin struct
			type FAB struct {
				Gateway string `json:"gateway"`
				Mtu     int    `json:"mtu"`
				Prefix  string `json:"prefix"`
				Vlan    int    `json:"vlan"`
				OOB     bool   `json:"out_of_band"`
			}

			//assign out of band values for api call.
			var OOB bool
			if subnetOOB.Selected() == 0 {
				OOB = true
			}
			if subnetOOB.Selected() == 1 {
				OOB = false
			}
			//convert the vlan entry from string to Int
			vlanInt, err := strconv.Atoi(subnetVlan.Text())
			//initialize FAS struct object
			FB := &FAB{}
			FB.OOB = OOB
			FB.Gateway = subnetGateway.Text()
			FB.Mtu = 1500
			FB.Prefix = subnetPrefix.Text()
			//FB.Services = svc
			FB.Vlan = vlanInt

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//make the patch request
			result := apiCallFB("PATCH", apiUrl+"/subnets?names="+subnetName.Text(), xAuthToken, FBData)
			initResult.SetText(string(result))
		}
	})

	subnetDeleteButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		err := validate.Var(subnetName.Text(), "required")
		if err != nil {
			initResult.SetText("Please provide the Subnet Name To Delete")
			passed = false
		}
		if passed == true {
			result := apiCallFB("DELETE", apiUrl+"/subnets?names="+subnetName.Text(), xAuthToken, nil)
			initResult.SetText(string(result))
		}
	})

	nicGetButton.OnClicked(func(*ui.Button) {
		result := apiCallFB("GET", apiUrl+"/network-interfaces", xAuthToken, nil)
		initResult.SetText(string(result))
	})

	nicPatchButton.OnClicked(func(*ui.Button) {
		initResult.SetText("processing vir0...")
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate
		err := validate.Var(virIP.Text(), "ip")
		if err != nil {
			initResult.SetText("Please provide a valid IP for VIR")
			passed = false
		}
		err1 := validate.Var(fm1AdminIP.Text(), "ip")
		if err1 != nil {
			initResult.SetText("Please provide a valid IP for FM1")
			passed = false
		}
		err2 := validate.Var(fm2AdminIP.Text(), "ip")
		if err2 != nil {
			initResult.SetText("Please provide a valid IP for FM2")
			passed = false
		}

		if passed == true {

			type FAB struct {
				Address string `json:"address"`
			}

			//initialize FAS struct object
			FB := &FAB{}
			FB.Address = virIP.Text()

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//initialize FAS struct object
			FB1 := &FAB{}
			FB1.Address = fm1AdminIP.Text()

			//marshal (json encode) the map into a json string
			FBData1, err := json.Marshal(FB1)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//initialize FAS struct object
			FB2 := &FAB{}
			FB2.Address = fm2AdminIP.Text()

			//marshal (json encode) the map into a json string
			FBData2, err := json.Marshal(FB2)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//variables for the PATCH request to apply NIC IP's fro VIR0, FM1, and FM2
			var r map[string]interface{}
			var vir0Status, fm1Status, fm2Status int
			var vir0Result, fm1Result, fm2Result string
			//poor mans async call.
			//start with VIR0 IP PATCH to Array, then check status.
			result0 := apiCallFB("PATCH", apiUrl+"/network-interfaces?names=vir0", xAuthToken, FBData)
			time.Sleep(time.Second * 3)
			vir0Status = statusCode
			if vir0Status == 200 {
				json.Unmarshal([]byte(result0), &r)
				vir0, d, o, err := jsonparser.Get(result0, "items", "[0]", "address")
				fmt.Print(d, o, err)
				vir0Result = "Successfully Applied the Vir0 Address: " + string(vir0)
			} else {
				vir0Result = "Failed. Something went wrong with the Vir0 Address. \n\nStatus code: " + string(vir0Status) + "\n\nresponse: " + string(result0)
			}
			//FM1 IP PATCH to Array, then check status.
			result1 := apiCallFB("PATCH", apiUrl+"/network-interfaces?names=fm1.admin0", xAuthToken, FBData1)
			time.Sleep(time.Second * 3)
			fm1Status = statusCode
			if fm1Status == 200 {
				json.Unmarshal([]byte(result1), &r)
				FM1, d, o, err := jsonparser.Get(result1, "items", "[0]", "address")
				fmt.Print(d, o, err)
				fm1Result = "\n\nSuccessfully Applied the FM1 Address: " + string(FM1)
			} else {
				fm1Result = "\n\nFailed. Something went wrong with the FM1 Address. \n\nStatus code: " + string(fm1Status) + "\n\nresponse: " + string(result1)
			}
			//FM2 IP PATCH to array then check status
			result2 := apiCallFB("PATCH", apiUrl+"/network-interfaces?names=fm2.admin0", xAuthToken, FBData2)
			time.Sleep(time.Second * 3)
			fm2Status = statusCode
			if fm2Status == 200 {
				json.Unmarshal([]byte(result2), &r)
				FM2, d, o, err := jsonparser.Get(result2, "items", "[0]", "address")
				fm2Result = "\n\nSuccessfully Applied the FM2 Address: " + string(FM2)
				fmt.Print(d, o, err)
			} else {
				fm2Result = "\n\nFailed. Something went wrong with the FM2 Addreess. \n\nStatus code: " + string(fm2Status) + "\n\nResponse: " + string(result2)
			}

			if vir0Status == 200 && fm1Status == 200 && fm2Status == 200 {
				if t, found := steps[5]; found {
					fmt.Println("step ", t, ", already completed not changing progress bar.")
				} else {
					steps[4] = "complete"
					progressCounter = progressCounter + 30
					prog.SetValue(progressCounter)
				}
				initResult.SetText(vir0Result + fm1Result + fm2Result + "\n\nPlease proceed to the SMTP Form.")
			} else {
				initResult.SetText(vir0Result + fm1Result + fm2Result)
			}
		}
	})

	smtpGetButton.OnClicked(func(*ui.Button) {
		result := apiCallFB("GET", apiUrl+"/smtp", xAuthToken, nil)
		initResult.SetText(string(result))
	})

	smtpPatchButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate
		//validate smtp domain is fqdn
		err1 := validate.Var(smtpSenderDomain.Text(), "fqdn")
		if err1 != nil {
			initResult.SetText("Please provide a sender domain as a fully qualified domain name.")
			passed = false
		}
		//if relay is set validate format
		if smtpRelayHost.Text() != "" {
			err := validate.Var(smtpRelayHost.Text(), "fqdn|ip|url")
			if err != nil {
				initResult.SetText("Please provide a valid IP, URL, or FQDN for the relay host")
				passed = false
			}
		}

		//if validation passes, continue
		if passed == true {

			type FAB struct {
				Relay  string `json:"relay_host"`
				Domain string `json:"sender_domain"`
			}

			//initialize FB struct object
			FB := &FAB{}
			FB.Relay = smtpRelayHost.Text()
			FB.Domain = smtpSenderDomain.Text()

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//make the rest call to the array api
			result := apiCallFB("PATCH", apiUrl+"/smtp", xAuthToken, FBData)
			//check for a successful Patch request
			if statusCode == 200 {
				//decode the response from a successfull call to a map interface//
				var r map[string]interface{}
				json.Unmarshal([]byte(result), &r)
				senderDomain, d, o, err := jsonparser.Get(result, "items", "[0]", "sender_domain")
				fmt.Print("jsonparser out: ", d, o, err)
				relay, d, o, err := jsonparser.Get(result, "items", "[0]", "relay_host")
				fmt.Print("jsonparser out: ", d, o, err)
				//print results to output field
				initResult.SetText("Success!\n\nApplied the following:\n\nSender Domain: " + string(senderDomain) + "\nRelay Host: " + string(relay) + "\n\nPlease proceed to the Support Form")
				//used for the progress bar
				if t, found := steps[6]; found {
					fmt.Println("step ", t, ", already completed not changing progress bar.")
				} else {
					steps[6] = "complete"
					progressCounter = progressCounter + 10
					prog.SetValue(progressCounter)
				}
			} else {
				initResult.SetText("Error: " + string(result))
			}
		}
	})

	supportGetButton.OnClicked(func(*ui.Button) {
		result := apiCallFB("GET", apiUrl+"/support", xAuthToken, nil)
		initResult.SetText(string(result))
	})

	supportPatchButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate
		if supportPhoneHome.Selected() == -1 {
			initResult.SetText("Please select true or false for Phone Home")
			passed = false
		}
		//set phonehome value for api call
		var phoneHome = false
		if supportPhoneHome.Selected() == 0 {
			phoneHome = true
		}
		if supportPhoneHome.Selected() == 1 {
			phoneHome = false
		}
		//if proxy is set validate format
		if supportProxy.Text() != "" {
			err := validate.Var(supportProxy.Text(), "fqdn|ip|url")
			if err != nil {
				initResult.SetText("Please provide a valid IP, URL or FQDN for the proxy server")
				passed = false
			}
		}

		if passed == true {

			type FAB struct {
				Phonehome bool   `json:"phonehome_enabled"`
				Proxy     string `json:"proxy"`
			}

			//initialize FAS struct object
			FB := &FAB{}
			FB.Phonehome = phoneHome
			FB.Proxy = supportProxy.Text()

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			result := apiCallFB("PATCH", apiUrl+"/support", xAuthToken, FBData)
			fmt.Print(string(result) + "\n")
			fmt.Print(statusCode)

			if statusCode == 200 {
				//decode the response from a successfull call to a map interface//
				var r map[string]interface{}
				json.Unmarshal([]byte(result), &r)
				phoneHomeEnabled, d, o, err := jsonparser.Get(result, "items", "[0]", "phonehome_enabled")
				fmt.Print("jsonparser out: ", d, o, err)
				proxy, d, o, err := jsonparser.Get(result, "items", "[0]", "proxy")
				fmt.Print("jsonparser out: ", d, o, err)
				//print results to output field
				initResult.SetText("Success!\n\nApplied the following:\n\nPhone Home Enabled: " + string(phoneHomeEnabled) + "\nProxy Server: " + string(proxy) + "\n\nPlease proceed to the Alerts Form")
				//used for the progress bar
				if t, found := steps[7]; found {
					fmt.Println("step ", t, ", already completed not changing progress bar.")
				} else {
					steps[7] = "complete"
					progressCounter = progressCounter + 10
					prog.SetValue(progressCounter)
				}
			} else {
				initResult.SetText("Error: " + string(result))
			}
		}
	})

	awGetButton.OnClicked(func(*ui.Button) {
		result := apiCallFB("GET", apiUrl+"/alert-watchers", xAuthToken, nil)
		initResult.SetText(string(result))
	})

	awPatchButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate
		err2 := validate.Var(awName.Text(), "email")
		if err2 != nil {
			initResult.SetText("Please provide an email address")
			passed = false
		}
		if awEnabled.Selected() == -1 {
			initResult.SetText("Please select true or false for enabled")
			passed = false
		}

		var awIsEnabled = ""
		if awEnabled.Selected() == 0 {
			awIsEnabled = "true"
		}
		if awEnabled.Selected() == 1 {
			awIsEnabled = "false"
		}

		if passed == true {

			type FAB struct {
				Enabled string `json:"enabled"`
			}

			//initialize FAS struct object
			FB := &FAB{}
			FB.Enabled = awIsEnabled

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//make the patch rest call to the array api
			result := apiCallFB("PATCH", apiUrl+"/alert-watchers?names="+awName.Text(), xAuthToken, FBData)
			//check if it was a success and provide the output.
			if statusCode == 200 {
				//decode the response from a successfull call to a map interface//
				var r map[string]interface{}
				json.Unmarshal([]byte(result), &r)
				email, d, o, err := jsonparser.Get(result, "items", "[0]", "name")
				fmt.Print("jsonparser out: ", d, o, err)
				enabled, d, o, err := jsonparser.Get(result, "items", "[0]", "enabled")
				fmt.Print("jsonparser out: ", d, o, err)
				//print results to output field
				initResult.SetText("Success!\n\nApplied the following:\n\nEmail Added: " + string(email) + "\nEnabled: " + string(enabled) + "\n\nPlease proceed to the Validation Form")
			} else {
				initResult.SetText("Error: " + string(result))
			}
		}
	})

	awDeleteButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate
		err := validate.Var(awName.Text(), "email")
		if err != nil {
			initResult.SetText("Please provide an email address")
			passed = false
		}
		if passed == true {
			result := apiCallFB("DELETE", apiUrl+"/alert-watchers?names="+awName.Text(), xAuthToken, nil)
			initResult.SetText(string(result))
		}
	})

	awPostButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true
		validate := validator.New()

		//validate alert watcher email
		err2 := validate.Var(awName.Text(), "email")
		if err2 != nil {
			initResult.SetText("Please provide an email address")
			passed = false
		}
		//validate and store if aw is enabled or disabled.
		var awIsEnabled = ""
		if awEnabled.Selected() > -1 {
			if awEnabled.Selected() == 0 {
				awIsEnabled = "true"
			}
			if awEnabled.Selected() == 1 {
				awIsEnabled = "false"
			}
		} else {
			initResult.SetText("Please select true or false for Enabled.")
			passed = false
		}

		// if all validation passes proceed
		if passed == true {

			type FAB struct {
				Enabled string `json:"enabled"`
			}
			//initialize FAS struct object
			FB := &FAB{}
			FB.Enabled = awIsEnabled

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//make the POST rest api call to the array
			result := apiCallFB("POST", apiUrl+"/alert-watchers?names="+awName.Text(), xAuthToken, FBData)
			//check if it was successful and present results.
			if statusCode == 200 {
				//decode the response from a successfull call to a map interface//
				var r map[string]interface{}
				json.Unmarshal([]byte(result), &r)
				email, d, o, err := jsonparser.Get(result, "items", "[0]", "name")
				fmt.Print("jsonparser out: ", d, o, err)
				enabled, d, o, err := jsonparser.Get(result, "items", "[0]", "enabled")
				fmt.Print("jsonparser out: ", d, o, err)
				//print results to output field
				initResult.SetText("Success!\n\nApplied the following:\n\nEmail Added: " + string(email) + "\nEnabled: " + string(enabled) + "\n\nPlease proceed to the Validation Form")
				//used for the progress bar
				if t, found := steps[8]; found {
					fmt.Println("step ", t, ", already completed not changing progress bar.")
				} else {
					steps[8] = "complete"
					progressCounter = progressCounter + 10
					prog.SetValue(progressCounter)
				}
			} else {
				initResult.SetText("Error: " + string(result))
			}
		}
	})

	finalGetButton.OnClicked(func(*ui.Button) {
		result := apiCallFB("GET", apiUrl+"/setup/validation", xAuthToken, nil)
		initResult.SetText(string(result))
	})

	finalPatchButton.OnClicked(func(*ui.Button) {
		//form validation object instantiation
		var passed bool = true

		//validate
		if finalSetupComplete.Selected() == -1 {
			initResult.SetText("Please select true or false")
			passed = false
		}

		var finalSetupCompleteIsComplete = ""
		if finalSetupComplete.Selected() == 0 {
			finalSetupCompleteIsComplete = "true"
		}
		if finalSetupComplete.Selected() == 1 {
			finalSetupCompleteIsComplete = "false"
		}

		if passed == true {

			type FAB struct {
				Complete string `json:"setup_completed"`
			}

			//initialize FAS struct object
			FB := &FAB{}
			FB.Complete = finalSetupCompleteIsComplete

			//marshal (json encode) the map into a json string
			FBData, err := json.Marshal(FB)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			result := apiCallFB("PATCH", apiUrl+"/setup/finalization", xAuthToken, FBData)
			initResult.SetText(string(result))
		}
	})
	//END Button Actions from forms//

	return hbox
}

//FUNCTION TO SETUP UI INTERFACE//
func setupUI() {
	mainwin = ui.NewWindow("Pure Storage Zero Touch Provisioner", 500, 300, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("ZTP Flash Array", initializeFATab())
	tab.SetMargined(0, true)

	tab.Append("ZTP Flash Blade", initializeFBTab())
	tab.SetMargined(0, true)

	mainwin.Show()
}

//MAIN//
func main() {
	ui.Main(setupUI)
}
