# iniset : A simple ini setter for config netdata based on env

## How to use it

export any number of env variables in this pattern :

`ND_{ANYTHING}=path/inifile.conf|section/key=value`

the `ND_` is default prefix and can change it by `-prefix=ANYOTHER_`
the path to the ini file is from `/etc/netdata` but it can changed with `-root=/abs/folder`
