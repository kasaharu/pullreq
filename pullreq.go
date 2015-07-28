package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	flag.Parse()
	fmt.Printf("------------------------------\n")
	if len(flag.Args()) <= 0 {
		fmt.Printf("invalid args \n")
		return
	}
	check_arg, err_for_matching := regexp.MatchString("^[0-9]+$", flag.Arg(0))
	if check_arg == false {
		fmt.Println(err_for_matching)
		fmt.Printf("Please input Ticket number : ex) hello 9999 \n")
		return
	} else {
		fmt.Printf("Create the PullRequest of No. %s \n", flag.Arg(0))
	}

	ticket_no := flag.Arg(0)
	template_file_name := "~/.gh_message_templates/pull_request.txt"
	temporary_file_name := "~/.gh_message_templates/pull_request_temp.txt"
	fmt.Printf("Create temporary format file.\n")

	sed_command := "sed -e \"s/ticket-no/" + ticket_no + "/g\" " + template_file_name + " >" + temporary_file_name
	result_sed, err_for_sed := exec.Command(os.Getenv("SHELL"), "-c", sed_command).Output()
	if err_for_sed != nil {
		fmt.Println(err_for_sed)
		return
	}
	fmt.Println(string(result_sed))

	fmt.Printf("Exec Pull Request.\n")
	result_hub_cmd, err_for_hub_cmd := exec.Command(os.Getenv("SHELL"), "-c", "hub pull-request --browse -F "+temporary_file_name+" -b kasaharu:master -h kasaharu:$(git rev-parse --abbrev-ref HEAD)").Output()
	if err_for_hub_cmd != nil {
		fmt.Println(err_for_hub_cmd)
		return
	}
	fmt.Println(string(result_hub_cmd))

	fmt.Printf("Delete temporary format file.\n")
	result_rm_temp, err_for_rm_temp := exec.Command(os.Getenv("SHELL"), "-c", "rm "+temporary_file_name).Output()
	if err_for_rm_temp != nil {
		fmt.Println(err_for_rm_temp)
		return
	}
	fmt.Println(string(result_rm_temp))
	fmt.Printf("------------------------------\n")
}
