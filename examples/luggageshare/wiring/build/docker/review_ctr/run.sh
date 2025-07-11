#!/bin/bash

WORKSPACE_NAME="review_ctr"
WORKSPACE_DIR=$(pwd)

usage() { 
	echo "Usage: $0 [-h]" 1>&2
	echo "  Required environment variables:"
	
	if [ -z "${REVIEW_DB_DIAL_ADDR+x}" ]; then
		echo "    REVIEW_DB_DIAL_ADDR (missing)"
	else
		echo "    REVIEW_DB_DIAL_ADDR=$REVIEW_DB_DIAL_ADDR"
	fi
	if [ -z "${REVIEW_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "    REVIEW_SERVICE_HTTP_BIND_ADDR (missing)"
	else
		echo "    REVIEW_SERVICE_HTTP_BIND_ADDR=$REVIEW_SERVICE_HTTP_BIND_ADDR"
	fi
		
	exit 1; 
}

while getopts "h" flag; do
	case $flag in
		*)
		usage
		;;
	esac
done


review_proc() {
	cd $WORKSPACE_DIR
	
	if [ -z "${REVIEW_DB_DIAL_ADDR+x}" ]; then
		if ! review_db_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${REVIEW_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		if ! review_service_http_bind_addr; then
			return $?
		fi
	fi

	run_review_proc() {
		
        cd review_proc
        ./review_proc --review_db.dial_addr=$REVIEW_DB_DIAL_ADDR --review_service.http.bind_addr=$REVIEW_SERVICE_HTTP_BIND_ADDR &
        REVIEW_PROC=$!
        return $?

	}

	if run_review_proc; then
		if [ -z "${REVIEW_PROC+x}" ]; then
			echo "${WORKSPACE_NAME} error starting review_proc: function review_proc did not set REVIEW_PROC"
			return 1
		else
			echo "${WORKSPACE_NAME} started review_proc"
			return 0
		fi
	else
		exitcode=$?
		echo "${WORKSPACE_NAME} aborting review_proc due to exitcode ${exitcode} from review_proc"
		return $exitcode
	fi
}


run_all() {
	echo "Running review_ctr"

	# Check that all necessary environment variables are set
	echo "Required environment variables:"
	missing_vars=0
	if [ -z "${REVIEW_DB_DIAL_ADDR+x}" ]; then
		echo "  REVIEW_DB_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  REVIEW_DB_DIAL_ADDR=$REVIEW_DB_DIAL_ADDR"
	fi
	
	if [ -z "${REVIEW_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "  REVIEW_SERVICE_HTTP_BIND_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  REVIEW_SERVICE_HTTP_BIND_ADDR=$REVIEW_SERVICE_HTTP_BIND_ADDR"
	fi
		

	if [ "$missing_vars" -gt 0 ]; then
		echo "Aborting due to missing environment variables"
		return 1
	fi

	review_proc
	
	wait
}

run_all
