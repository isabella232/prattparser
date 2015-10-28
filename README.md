prattparser
=====

A sample implementation of a parser based on Pratt algorithm for the xoom blog [Parsing made easy by Pratt algorithm - An Implementation in Go](http://dev-blog.xoom.com/2015/10/26/)

Install
=====

     go get github.com/xoom/prattparser

which installs into $GOPATH/bin/prattparser

Usage
=====

     prattparser -input " 1 + var1 > 3" -variables="var1=5, var3=4"

which will print "evaluated result is  true".
