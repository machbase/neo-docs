# Machbase Neo Documents

This repository is for the **Machbase Neo** documentation published at [machbase.github.io/machbase](https://machbase.github.io/machbase)

# Writer's environment

## Ubuntu

1. `apt-get install ruby-full build-essential zlib1g-dev`
2. add PATH
    
    ```
    export GEM_HOME="$HOME/gems"
    export PATH="$HOME/gems/bin:$PATH"
    ```

3. `gem install jekyll bundler`
4. `bundle install`
5. `bundle exec jekyll serve`

## macOS

1. `brew install chruby ruby-install xz`
2. `ruby-install ruby 3.1.3`
3. .zshrc
```
echo "source $(brew --prefix)/opt/chruby/share/chruby/chruby.sh" >> ~/.zshrc
echo "source $(brew --prefix)/opt/chruby/share/chruby/auto.sh" >> ~/.zshrc
echo "chruby ruby-3.1.3" >> ~/.zshrc # run 'chruby' to see actual version
```

4. `gem install jekyll bundler`
5. `bundle install`
6. `bundle exec jekyll serve`

