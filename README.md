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

4. Call echo method default grpc endpoint is `localhost:9876`

5. You'll see logs on your terminal
