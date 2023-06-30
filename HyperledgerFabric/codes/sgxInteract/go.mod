module sgxInteract

go 1.20


require myrsa v0.0.0-00010101000000-000000000000
replace myrsa => ../toolPackages/myrsa
require msgs v0.0.0-00010101000000-000000000000
replace msgs => ../toolPackages/msgs
require sharedConfigs v0.0.0-00010101000000-000000000000
replace sharedConfigs => ../toolPackages/sharedConfigs
