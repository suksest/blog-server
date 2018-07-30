# blog-server

## Description 
Restful API Web server for simple blog using Echo framework, Gorm, and Postgresql. 

## Installation
Clone the project:
```
  git clone https://github.com/ridwanfathin/blog-server
 ```
 open the directory:
 ```
  cd echo-server
 ```
 
 define `GOPATH`:
 ```
  export GOPATH=[your project path]
 ```
 
 build binaries: 
 ``` 
 go build main
 
 ```
 
 run the server: 
 
 ```
 ./main (linux & Mac)
 main.exe (windows)

 
 ```
 
 Note: if you want to use anything other the `master` aka the last publish part, switch to
 its branch and build binaries.
 
 list all branches:
 ```
  git branch -a
 ```
 checkout to the branch you want and build/run like before:
 ```
  git checkout part_1_hello_world
 
  go install main
 
  bin/main
 ```
  
