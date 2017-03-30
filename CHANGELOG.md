### Recent Updates

**22nd March 2017**: 
- More reorganisation: REST package seperated out from core and into its own package.  
- First draft of GraphQL package completed with simple list viewing capability.  
- Reorganisation of initialisation procedure to fix bug which created config.json files when testing.
- Custom config files can be used on any command with `-c` flag (don't add the extension) 

**16th March 2017**: I've now completed the first large refactor and updated master.  The Gin router has been switched for Chi, which is is compatible with default Go handlers, making life a little easier. Test coverage is coming along nicely. Functionality has been reorganised into 'packages' (which are just Go packages), allowing for much easier extensibility.  *Core, Auth and Email* packages are included in the standard build - everything else (HTML server, image resizer and admin-panel server) has been stripped out and will be rewritten and uploaded as seperate packages that can be added to a [custom build](https://github.com/jpincas/ghost-custom-server).  We're rethinking our strategy on the admin panel side of things, as the Polymer app was proving hard to install for newcomers and was causing some build problems when trying to reference external HTML imports.  We might go with Elm.  I'm also considering the possiblity of ditching the HTML server in favour of a static site integrator powered by Hugo.  Please get in touch if you'd like to help out!