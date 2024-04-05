#!/bin/sh

exec tclsh "$0" ${1+$@}

puts "Heyy this worked!"
