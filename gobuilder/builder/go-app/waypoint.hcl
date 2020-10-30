project = "guides"

app "example" {
    build {
        use "gobuilder" {
            output_name = "simple-app"
            source = "./"
        }
    }
    deploy {
        use "godeploy" {
            directory = "something"
        }
    }
}