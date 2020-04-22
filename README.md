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