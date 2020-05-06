[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![CLA assistant](https://cla-assistant.io/readme/badge/PureStorage-OpenConnect/zero-touch-provisioner)](https://cla-assistant.io/PureStorage-OpenConnect/zero-touch-provisioner)

# FlashArray and FlashBlade Zero Touch Provisioner

Pure Storage FlashArray and FlashBlade Zero Touch Provisioner. This lightweight application allows engineers (or customers) to remotely initialize a new FA or FB using this cross-platform utility.

## Getting Started

Compiled binaries for Mac, Linux, and Windows are located in the Compiled directory.

Download the application binary relevant to your desktop OS and execute it.  Fill out the fields and your off.

The single console application has two tabs: FlashArray and FlashBlade.

#### FlashArray

![Blank Tab for FlashArray](assets/FA.PNG)

The ZTP FlashArray tab is comprised of 3 columns.  Columnm 1 is the general configurations sections, Column 2 is the IP Addressing section and Column 3 is the Action and Output section.

To begin you will enter the Temp DHCP IP of FlashArray in the "DHCP IP of Array" field in column 3.  You can then Query the Array to make sure you have connectivity.

Next you will fill out the remainder of the fileds in Column 1 and Column 2 according to the specifics of the Array and environment.

When ready to initialize, press the Initialize Array button in Column 3 and view the results in the output window below. 

#### FlashBlade

![Blank FlashBlade Tab](assets/FB.PNG)

Select the ZTP FlashBlade Tab.

The ZTP FlashBlade tab is comprised of 3 columns as well.  Column 1 provides the 9 Steps necessary to configure a FlashBlade in order of operation.  Each step contains a button that will display the relevant form and controls for the associated step.  Column 2 is the Form fields and controls to query and apply required configurations relevant to each of the 9 steps.  Column 3 provides the output of each action.

In order to authenticate to a FlashBlade using the ZTP tool, you will need the DHCP IP of the FlashBlade.

To begin you will can either query the FlashBlade to obtain the API versions and auto generate the API url for the "Array API URL" field or you can manually enter the API url in the "Array API URL" field.  If you manually enter the API URL you must include the version and have the format as http or ```https://<fqdn or ip of array>/api/<api version>```

![Query the Array builds the API URL with the latest support version](assets/FBQuery.PNG)

When you have the API url in place, clicking the "Create Session" button should result in the API URL and x-auth-token labels in column 3 displaying the resulting values.  If you do not see a value listed in the x-auth-token label field then the session was not created.

![Create Session logs into the array using the api-token and gets a session token called x-auth-token](assets/FBCreateSession.PNG)

When the above is successful you are ready to move to Step 2 - Step 9 in succession.

### Prerequisites

Windows 10 or higher, Linux with UI, or MAC OS.

```
For Windows, download and launch the appropriate .exe, fill out the form and go.
```

```
For Linux, (desktop GUI is required) download and launch appropriate ZTP binary, set it to executable "chmod +x", execute it and fill out the form.
```

```
For MacOS, download the appropriate ZTP binary, open terminal, set the binary to executable "chmod +x", execute it and fill out the form.
```

## Video Tutorial

[![Video Tutorial](assets/ZTP_HowTo_First_Frame.png)](https://youtu.be/LN13cpuA3Jc)

## Screen Shots

FlashArray

![Blank Tab for FlashArray](assets/FA.PNG)

FlashBlade

![Blank Tab for FlashBlade](assets/FB.PNG)

## Built With

* [UI](https://github.com/andlabs/ui/) - The GUI framework used
* [Validator v9](https://gopkg.in/go-playground/validator.v9) - form validation
* [GOLang](https://golang.org/) - Core
* [Crossbuild](https://github.com/magJ/go-ui-crossbuild) - Compiled with Go UI Crossbuild


## Authors

* **Brandon Showers* - - [GIBLAB](https://gitlab.com/btshowers)  OR  [About Me](https://bshowers.com)


## Acknowledgments

* Thanks Simon Dodsley for the hackathon to get me to learn GOLang!

