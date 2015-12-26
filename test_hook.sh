#!/bin/bash

read -d '' PAYLOAD <<"EOF"
{
  "ref": "refs/heads/master",
  "after": "6ead527eb6a7a495852706ca3f6e3d715e0ef141",
  "deleted": false,
  "repository": {
    "full_name": "user/repo"
  }
}
EOF

curl --data-binary "$PAYLOAD" -H "X-Hub-Signature: sha1=a45a91e911575ad2933640a5ec0bc36df8d479f9" gohub_test.docker