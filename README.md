# soilgaze
Soilgaze is a simple Go application that recovers intel from OSINT sources such as Shodan, Binaryedge and Censys.


## OSINT Resources

The resources below are either already integrated within the application or planned to be integrated in the future.

| Resources     | Current Status    |
| ------------- | :---------------: |
| Shodan        | COMPLETED         |
| Binaryedge    | COMPLETED         |
| Censys        | COMPLETED         |
| Zoomeye       | IN-PROGRESS       |
| Onyphe        | COMPLETED         |
| Spyse         | IN-PROGRESS       |

## Flags

| Flags         | Functions                                              |
| ------------- | ------------------------------------------------------ |
| host-file     | The location of the file that contains a list of hosts |
| osint-list    | Comma-separated string that contains OSINT resources   |
| out-file      | Destination to save JSON output of enumeration results |

## Current Situation
* Spyse has multiple issues that prevents me from implementing it:
  * There are no endpoints for checking remaining query rights. Even the web application does not show such information.
  * API is returning unstable responses. Currently getting 401 even with a brand new account.
* Zoomeye also has an issue. It requires a phone number to register which I will not give, naturally.