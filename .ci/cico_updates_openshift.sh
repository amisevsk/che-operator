#!/bin/bash
#
# Copyright (c) 2020 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   Red Hat, Inc. - initial API and implementation

set -ex

#Stop execution on any error
trap "catchFinish" EXIT SIGINT

# Catch_Finish is executed after finish script.
catchFinish() {
  exit $result
}

init() {
  export SCRIPT=$(readlink -f "$0")
  export SCRIPT_DIR=$(dirname "$SCRIPT")

  if [[ ${WORKSPACE} ]] && [[ -d ${WORKSPACE} ]]; then
    export OPERATOR_REPO=${WORKSPACE};
  else
    export OPERATOR_REPO=$(dirname "$SCRIPT_DIR");
  fi

  export RAM_MEMORY=8192
  export PLATFORM="openshift"
  export NAMESPACE="che"
  export CHANNEL="stable"
}

waitCheUpdateInstall() {
  export packageName=eclipse-che-preview-${PLATFORM}
  export platformPath=${OPERATOR_REPO}/olm/${packageName}
  export packageFolderPath="${platformPath}/deploy/olm-catalog/${packageName}"
  export packageFilePath="${packageFolderPath}/${packageName}.package.yaml"

  export lastCSV=$(yq -r ".channels[] | select(.name == \"${CHANNEL}\") | .currentCSV" "${packageFilePath}")
  export lastPackageVersion=$(echo "${lastCSV}" | sed -e "s/${packageName}.v//")

  echo -e "\u001b[34m Check installation last version che-operator...$lastPackageVersion \u001b[0m"

  export n=0

  while [ $n -le 360 ]
  do
    cheVersion=$(kubectl get checluster/eclipse-che -n "${NAMESPACE}" -o jsonpath={.status.cheVersion})
    if [ "${cheVersion}" == $lastPackageVersion ]
    then
      echo -e "\u001b[32m Installed latest version che-operator: ${lastCSV} \u001b[0m"
      break
    fi
    sleep 3
    n=$(( n+1 ))
  done

  if [ $n -gt 360 ]
  then
    echo "Latest version install for Eclipse che failed."
    exit 1
  fi
}

testUpdates() {
  "${OPERATOR_REPO}"/olm/testUpdate.sh ${PLATFORM} ${CHANNEL} ${NAMESPACE}
  printInfo "Successfully installed Eclipse Che previous version."

  getCheAcessToken
  chectl workspace:create --devfile=$OPERATOR_REPO/.ci/util/devfile-test.yaml

  waitCheUpdateInstall
  getCheAcessToken

  sleep 120
  local cheVersion=$(kubectl get checluster/eclipse-che -n "${NAMESPACE}" -o jsonpath={.status.cheVersion})

  echo "[INFO] Successfully installed Eclipse Che: ${cheVersion}"

  # Wait to start an workspace and print success message
  waitWorkspaceStart
  echo "[INFO] Successfully started an workspace on Eclipse Che: ${cheVersion}"
}

init
source "${OPERATOR_REPO}"/.ci/util/ci_common.sh
testUpdates
