data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./domain/models",  
    "--dialect", "sqlite"   
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "sqlite://data.db"  
  migration {
    dir = "file://migrations"  # Directory for storing migration files
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"  # Format for migration differences
    }
  }
}
