package main

func callSorlInstallTomcat(ss *SorlSSH, cmd string, allProp *Property) {

	if cmd != "tomcat" {
		ss.sshPrint("error: package name is not tomcat", allProp)
		return
	}

	tomcatUrl := (*allProp)["tomcat.url"]
	tomcatFile := (*allProp)["tomcat.file"]
	tomcatInstallPath := (*allProp)["tomcat.install.path"]
	tomcatVersion := (*allProp)["tomcat.version"]
	tomcatUserId := (*allProp)["tomcat.user.id"]
	tomcatUserGroupId := (*allProp)["tomcat.user.group.id"]
	tomcatListenPort := (*allProp)["tomcat.listen.port"]
	//tomcatAddUser := (*allProp)["tomcat.add.user"]

	strInfo := "\n"
	strInfo += "\nInstalling Tomcat "
	strInfo += "\nDetails:"
	strInfo += "\n\tVersion: " + tomcatVersion
	strInfo += "\n\tFile: " + tomcatFile
	strInfo += "\n\tPath: " + tomcatInstallPath
	strInfo += "\n\tUser: " + tomcatUserId
	strInfo += "\n\tGroup: " + tomcatUserGroupId
	strInfo += "\n\tPort: " + tomcatListenPort
	strInfo += "\n"
	ss.sshPrint(strInfo, allProp)

	if tomcatUrl != "" {
		ss.sshPrint("\ndownloading file '"+tomcatFile+"' to /tmp", allProp)
		prevEcho := (*allProp)["sr:echo"]
		(*allProp)["sr:echo"] = "on"

		(*allProp)["_if.prompt.req"] = "true"

		ss.runShellCmd("id -u " + tomcatUserId)
		callSorlOrchWait(ss, (*allProp)["_wait.prev.cmd"], allProp)
		callSorlOrchStatus(ss, "status", allProp)

		if (*allProp)["_exit.code"] != "0" {
			ss.runShellCmd("groupadd " + tomcatUserGroupId)
			callSorlOrchWait(ss, (*allProp)["_wait.prev.cmd"], allProp)

			ss.runShellCmd("useradd -m -g " + tomcatUserGroupId + " -d " + tomcatInstallPath + " " + tomcatUserId)
			callSorlOrchWait(ss, (*allProp)["_wait.prev.cmd"], allProp)

			ss.runShellCmd("usermod -a -G " + tomcatUserGroupId + " " + tomcatUserId)
			callSorlOrchWait(ss, (*allProp)["_wait.prev.cmd"], allProp)
		}

		ss.runShellCmd("cd /tmp && wget " + tomcatUrl)
		callSorlOrchWait(ss, (*allProp)["_wait.prev.cmd"], allProp)

		ss.runShellCmd("mkdir -p " + tomcatInstallPath)
		callSorlOrchWait(ss, (*allProp)["_wait.prev.cmd"], allProp)

		ss.runShellCmd("cd " + tomcatInstallPath)
		callSorlOrchWait(ss, (*allProp)["_wait.prev.cmd"], allProp)

		ss.runShellCmd("tar zxvf /tmp/" + tomcatFile)
		callSorlOrchWait(ss, (*allProp)["_wait.prev.cmd"], allProp)

		ss.runShellCmd("chown -R " + tomcatUserId + ":" + tomcatUserGroupId + " " + tomcatInstallPath)
		callSorlOrchWait(ss, (*allProp)["_wait.prev.cmd"], allProp)

		(*allProp)["sr:echo"] = prevEcho
	}

}
