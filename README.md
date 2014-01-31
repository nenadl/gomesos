gomesos
=======

A fork and reorganization of the https://github.com/mesosphere/mesos-go repository so that it's buildable via the go build system.

To install this project just do:

    go get github.com/nenadl/gomesos
  
To install the example framework and executor do

    go install github.com/nenadl/gomesos/...
  
On OSX with CLANG you should use GCC

    CC=gcc-4.8 CXX=g++-4.8 go get github.com/nenadl/gomesos
    
To include in your go projects do:

    import "github.com/nenadl/gomesos"
