project = "guides"

app "example" {
    build {
        use "gobuilder" {
            output_name = "lambda"
            source = "./"
            arch {
		amd64 = true 
            	linux = true
	    }
        }
    }
    deploy {
        use "godeploy" {
        }
    }
}
