#!/usr/bin/perl

sub print_parameters {
    my @params = @_;
    print "COUNT PARAMS: " . scalar(@params) . "\n";
    print "PARAMS:\n";
    foreach my $param (@params) {
        print "$param\n";
    }
}

print "PERL\n";
print_parameters(@ARGV);
