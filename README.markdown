a fun little fact database.

Written in Go and backed by Riak.

This previously existed (as Django+PostgreSQL) as
[faktum](https://github.com/thraxil/faktum/).

The idea is to give me a simple database for capturing little
facts/trivia I learn along with their sources, organized by tags.

This incarnation is largely and excuse to do some non-trivial data
modelling with Riak and to get more comfortable writing full webapps
with Go. It is what I use day to day though and can be seen in use at

    http://facts.thraxil.org/
