# Synopsis: Run tests
task test {
    Set-Location -Path 'app'
    $Status = Start-Process -FilePath 'go' -ArgumentList 'test' -NoNewWindow -PassThru -Wait
    Assert($Status.ExitCode -eq 0) 'The test command failed'
}

# Synopsis: Remove generated files
task clean {
    Set-Location -Path 'app'
    $Status = Start-Process -FilePath 'go' -ArgumentList 'clean' -NoNewWindow -PassThru -Wait 
    Assert($Status.ExitCode -eq 0) 'The "clean" command failed'
}

# Synopsis: Generate documentation
task doc {
    Set-Location -Path 'app'
    $Status = Start-Process -FilePath 'go' -ArgumentList 'doc -all' -RedirectStandardOutput 'readme.txt' -NoNewWindow -PassThru -Wait 
    Assert($Status.ExitCode -eq 0) 'The "go doc" command failed'
}

task . test, clean, doc