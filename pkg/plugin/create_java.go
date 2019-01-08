package plugin

import (
	"fmt"
	"github.com/softleader/slctl/pkg/strcase"
	"path/filepath"
	"strings"
)

const javaMain = `package tw.com.softleader.slctl.plugin;

import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import java.util.stream.Stream;

@SpringBootApplication
public class {{.Name|title|camel}}Application implements CommandLineRunner {

  public static void main(String[] args) {
    SpringApplication.run({{.Name|title|camel}}Application.class, args);
  }

  @Override
  public void run(String... args) throws Exception {
    System.getenv()
        .entrySet()
        .stream()
        .filter(e -> e.getKey().startsWith("SL_"))
        .forEach(System.out::println);
  }
}

`

const javaProperties = `spring.main.banner-mode=off
logging.level.root=error
`

const javaPom = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
	<modelVersion>4.0.0</modelVersion>

	<groupId>tw.com.softleader</groupId>
	<artifactId>{{.Name}}</artifactId>
	<version>{{.Version}}</version>
	<packaging>jar</packaging>

	<name>{{.Name}}</name>
	<description>{{.Description}}</description>

	<parent>
		<groupId>org.springframework.boot</groupId>
		<artifactId>spring-boot-starter-parent</artifactId>
		<version>2.1.0.RELEASE</version>
		<relativePath/> <!-- lookup parent from repository -->
	</parent>

	<properties>
		<project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
		<project.reporting.outputEncoding>UTF-8</project.reporting.outputEncoding>
		<java.version>1.8</java.version>
	</properties>

	<dependencies>
		<dependency>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter</artifactId>
		</dependency>

		<dependency>
			<groupId>org.springframework.boot</groupId>
			<artifactId>spring-boot-starter-test</artifactId>
			<scope>test</scope>
		</dependency>
	</dependencies>

	<build>
		<finalName>{{.Name}}</finalName>
		<plugins>
			<plugin>
				<groupId>org.springframework.boot</groupId>
				<artifactId>spring-boot-maven-plugin</artifactId>
			</plugin>
		</plugins>
	</build>
</project>
`

const javaMakefile = `SL_HOME ?= $(shell slctl home)
SL_PLUGIN_DIR ?= $(SL_HOME)/plugins/{{.Name}}/
METADATA := metadata.yaml
HAS_MAVEN := $(shell command -v mvn;)
HAS_JDK := $(shell command -v javac;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' $(METADATA))
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/target
BINARY := {{.Name}}

.PHONY: install
install: bootstrap build
	mkdir -p $(SL_PLUGIN_DIR)
	cp $(BUILD)/$(BINARY).jar $(SL_PLUGIN_DIR)
	cp $(METADATA) $(SL_PLUGIN_DIR)

.PHONY: build
build: clean bootstrap
	mvn clean package

.PHONY: dist
dist: bootstrap build
	mkdir -p $(DIST)
	sed -E 's/(version: )"(.+)"/\1"$(VERSION)"/g' $(METADATA) > $(BUILD)/$(METADATA)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-$(VERSION).tgz $(BINARY).jar $(METADATA)

.PHONY: bootstrap
bootstrap:
ifndef HAS_JDK
	$(error You must install JDK)
endif
ifndef HAS_MAVEN
	$(error You must install Maven)
endif

.PHONY: clean
clean:
	rm -rf _*
	mvn clean
`

type java struct{}

func (c java) exec(plugin *Metadata) Commands {
	command := fmt.Sprintf("java -jar $SL_PLUGIN_DIR/%s.jar", plugin.Name)
	return Commands{
		Command: command,
		Platform: []Platform{
			{Os: "darwin", Command: command},
			{Os: "windows", Command: command},
		},
	}
}

func (c java) hook(plugin *Metadata) Commands {
	return Commands{
		Command: "echo hello " + plugin.Name,
	}
}

func (c java) files(plugin *Metadata, pdir string) []file {
	main := filepath.Join(pdir, "src", "main")
	return []file{
		tpl{
			path:     filepath.Join(main, "java", "tw", "com", "softleader", "slctl", "plugin", strcase.ToCamel(strings.Title(plugin.Name))+"Application.java"),
			in:       plugin,
			template: javaMain,
		},
		tpl{
			path:     filepath.Join(main, "resources", "application.properties"),
			in:       plugin,
			template: javaProperties,
		},
		tpl{
			path:     filepath.Join(pdir, "pom.xml"),
			in:       plugin,
			template: javaPom,
		},
		tpl{
			path:     filepath.Join(pdir, "Makefile"),
			in:       plugin,
			template: javaMakefile,
		},
	}
}
