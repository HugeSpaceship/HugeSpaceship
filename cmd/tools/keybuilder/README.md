# Key Builder

This tool takes raw ECDSA keys and turns them into a PEM file for use in HugeSpaceship, 
this makes them easier to move around, and makes adding custom keys a breeze.

## Usage
```
keybuilder -x <ECDSA X Point> -y <ECDSA Y Point> -curve secp192r1 -owner "HugeSpaceship" -desc "HugeSpaceship NpTicket signing key" -signatory 12345678
```

