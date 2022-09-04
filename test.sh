  set -e

  echo 'mode: count' > profile.cov

  for dir in */ ; do
          cd "$dir"
          if [ ! -f go.mod ]; then
              cd ..
              continue
          fi
          go test -short -covermode=count -coverprofile=./profile.tmp  ./...
          cd ..
          if [ -f $dir/profile.tmp ]
          then
              cat $dir/profile.tmp | tail -n +2 >> profile.cov
              rm $dir/profile.tmp
          fi
  done

  go tool cover -func profile.cov