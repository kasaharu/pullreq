package main

import (
    "bufio"
    "flag"
    "fmt"
    "io/ioutil"
    "gopkg.in/yaml.v2"
    "os"
    "os/exec"
    "regexp"
)

func main() {
    flag.Parse()

    config := parse_setting_file()
    option := check_args(flag.Args())

    exec_cmd(option, config)
}

func check_args(args []string) int {
    // 引数があるかどうか？
    if len(flag.Args()) <= 0 {
        fmt.Printf("invalid args \n")
        os.Exit(0)
    }

    // setting コマンドかどうか？
    is_set_cmd, err_is_setting := regexp.MatchString("setting", flag.Arg(0))
    if is_set_cmd == true {
        return 1
    } else {
        if err_is_setting != nil {
            fmt.Println(err_is_setting)
        }
    }

    // チケット番号かどうか？
    check_arg, err_for_matching := regexp.MatchString("^[0-9]+$", flag.Arg(0))
    if check_arg == true {
        fmt.Printf("Create the PullRequest of No. %s \n", flag.Arg(0))
        return 2
    } else {
        fmt.Println(err_for_matching)
        fmt.Printf("Please input Ticket number : ex) hello 9999 \n")
        os.Exit(0)
    }
    return 0
}

func exec_cmd(option int, config map[interface{}]interface{}) {
    base_user_name    := config["base_user_name"].(string)
    base_branch_name  := config["base_branch_name"].(string)
    compare_user_name := config["compare_user_name"].(string)

    switch {
    case option == 1 :
        cmd_type := select_set_cmd_type();
        fmt.Println();
        is_show_conf, _ := regexp.MatchString("1", cmd_type)
        if is_show_conf == true {
            fmt.Println("base_user_name    : ", base_user_name)
            fmt.Println("base_branch_name  : ", base_branch_name)
            fmt.Println("compare_user_name : ", compare_user_name)
        }
        is_set_tempfile, _ := regexp.MatchString("2", cmd_type)
        if is_set_tempfile == true {
            fmt.Printf("Setting template file. \n")
            result_copy_template, err_copy_template := exec.Command(os.Getenv("SHELL"), "-c", "cp -r $GOPATH/src/github.com/kasaharu/pullreq/.gh_message_templates ~/").Output()
            if err_copy_template != nil {
                fmt.Println("Fail to copy template file.\n")
                return
            } else {
                fmt.Println(string(result_copy_template))
                fmt.Printf("Success. \n")
                os.Exit(0)
            }
        }
        return
    case option == 2 :
        ticket_no := flag.Arg(0)
        template_file_name  := "~/.gh_message_templates/pull_request.txt"
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
        result_hub_cmd, err_for_hub_cmd := exec.Command(os.Getenv("SHELL"), "-c", "hub pull-request --browse -F "+temporary_file_name+" -b "+base_user_name+":"+base_branch_name+" -h "+compare_user_name+":$(git rev-parse --abbrev-ref HEAD)").Output()
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
            return
    }
    return
}

func parse_setting_file() map[interface{}]interface{} {
    setting_file := os.Getenv("GOPATH") + "/src/github.com/kasaharu/pullreq/setting/config.yml"
    buf, err := ioutil.ReadFile(setting_file)
    if err != nil {
        fmt.Println(err)
        os.Exit(0)
    }
    m := make(map[interface{}]interface{})
    err = yaml.Unmarshal(buf, &m)
    if err != nil {
        fmt.Println(err)
        os.Exit(0)
    }

    return m
}

func select_set_cmd_type() string {
    reader := bufio.NewReader(os.Stdin);
    fmt.Printf("Select the number of set command option.\n")
    fmt.Printf("  - 1 : Show the params of config.\n")
    fmt.Printf("  - 2 : Set the template file.\n")
    input, _ := reader.ReadString('\n');
    return input;
}
