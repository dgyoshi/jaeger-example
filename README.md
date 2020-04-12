# prerequisite
## OSX  
- direnv https://github.com/direnv/direnv
install direnv  
```
brew install direnv
```

Add following lines to .zshrc
```
export EDITOR=< your editor >
eval "$(direnv hook zsh)"
```

# try out
1. install bloomrpc https://github.com/uw-labs/bloomrpc
```
brew cask install bloomrpc
```

2. run docker compose
```
docker-compose up
```

3. Import echo.proto `proto/echo/echo.proto`

4. Call echo method default grpc endpoint is `localhost:9876`, you'll see trace_id in its response

5. Open your browser and access to http://localhost:16686/
6. Search and trace with the trace_id
