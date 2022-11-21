# network-monitor
Utility to measure Network Connectivity and Downtime.  Inspired by a need to measure and document ISP real world performance to inform provider choice. Uses email notifications to alert you of downtime resolutions and provide comprehensive log files which can be used to track long term network reliability.

## Setup
To build a local image:

`docker build -t network-monitor .`

To run:

`docker-compose up -d`

To watch logs in realtime:

`docker logs -f network-monitor`