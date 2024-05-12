# useful_things
This repo contains things that are useful exclusively for me for work, study and comfort


For use VirusTotal tools need build docker image first  
You cat do it from virustotal/Dockerfile   use  *docker build --no-cache --progress=plain --secret id=apikey,src=apikey.txt -t image_name:tag* command 
For run use *docker run -ti --rm image_name:tag*  command  By default you can see --help   
docker run -ti --rm  virustotal:latest 
A command-line tool for interacting with VirusTotal.

Usage:
  vt [command]

Available Commands:
  analysis       Get a file or URL analysis
  collection     Get information about collections
  completion     Output shell completion code for the specified shell (bash or zsh)
  domain         Get information about Internet domains
  download       Download files
  file           Get information about files
  group          Get information about VirusTotal groups
  help           Help about any command
  hunting        Manage malware hunting rules and notifications
  init           Initialize or re-initialize vt command-line tool
  iocstream      Manage IoC Stream notifications
  ip             Get information about IP addresses
  meta           Returns metadata about VirusTotal
  monitor        Manage your monitor account
  monitorpartner Manage your monitor partner account
  retrohunt      Manage retrohunt jobs
  scan           Scan files or URLs
  search         Search for files in VirusTotal Intelligence
  url            Get information about URLs
  user           Get information about VirusTotal users
  version        Show version number

Flags:
  -k, --apikey string   API key
      --format string   Output format (yaml/json/csv) (default "yaml")
  -h, --help            help for vt
      --proxy string    HTTP proxy
  -s, --silent          Silent or quiet mode. Do not show progress meter
  -v, --verbose         verbose output

Use "vt [command] --help" for more information about a command.
