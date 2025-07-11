#!/bin/bash

WORKSPACE_NAME="frontend_ctr"
WORKSPACE_DIR=$(pwd)

usage() { 
	echo "Usage: $0 [-h]" 1>&2
	echo "  Required environment variables:"
	
	if [ -z "${FRONTEND_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "    FRONTEND_SERVICE_HTTP_BIND_ADDR (missing)"
	else
		echo "    FRONTEND_SERVICE_HTTP_BIND_ADDR=$FRONTEND_SERVICE_HTTP_BIND_ADDR"
	fi
	if [ -z "${LUGGAGE_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "    LUGGAGE_SERVICE_HTTP_DIAL_ADDR (missing)"
	else
		echo "    LUGGAGE_SERVICE_HTTP_DIAL_ADDR=$LUGGAGE_SERVICE_HTTP_DIAL_ADDR"
	fi
	if [ -z "${RESERV_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "    RESERV_SERVICE_HTTP_DIAL_ADDR (missing)"
	else
		echo "    RESERV_SERVICE_HTTP_DIAL_ADDR=$RESERV_SERVICE_HTTP_DIAL_ADDR"
	fi
	if [ -z "${REVIEW_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "    REVIEW_SERVICE_HTTP_DIAL_ADDR (missing)"
	else
		echo "    REVIEW_SERVICE_HTTP_DIAL_ADDR=$REVIEW_SERVICE_HTTP_DIAL_ADDR"
	fi
	if [ -z "${SEARCH_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "    SEARCH_SERVICE_HTTP_DIAL_ADDR (missing)"
	else
		echo "    SEARCH_SERVICE_HTTP_DIAL_ADDR=$SEARCH_SERVICE_HTTP_DIAL_ADDR"
	fi
	if [ -z "${USER_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "    USER_SERVICE_HTTP_DIAL_ADDR (missing)"
	else
		echo "    USER_SERVICE_HTTP_DIAL_ADDR=$USER_SERVICE_HTTP_DIAL_ADDR"
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


frontend_proc() {
	cd $WORKSPACE_DIR
	
	if [ -z "${SEARCH_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		if ! search_service_http_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${RESERV_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		if ! reserv_service_http_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${REVIEW_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		if ! review_service_http_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${USER_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		if ! user_service_http_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${LUGGAGE_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		if ! luggage_service_http_dial_addr; then
			return $?
		fi
	fi

	if [ -z "${FRONTEND_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		if ! frontend_service_http_bind_addr; then
			return $?
		fi
	fi

	run_frontend_proc() {
		
        cd frontend_proc
        ./frontend_proc --search_service.http.dial_addr=$SEARCH_SERVICE_HTTP_DIAL_ADDR --reserv_service.http.dial_addr=$RESERV_SERVICE_HTTP_DIAL_ADDR --review_service.http.dial_addr=$REVIEW_SERVICE_HTTP_DIAL_ADDR --user_service.http.dial_addr=$USER_SERVICE_HTTP_DIAL_ADDR --luggage_service.http.dial_addr=$LUGGAGE_SERVICE_HTTP_DIAL_ADDR --frontend_service.http.bind_addr=$FRONTEND_SERVICE_HTTP_BIND_ADDR &
        FRONTEND_PROC=$!
        return $?

	}

	if run_frontend_proc; then
		if [ -z "${FRONTEND_PROC+x}" ]; then
			echo "${WORKSPACE_NAME} error starting frontend_proc: function frontend_proc did not set FRONTEND_PROC"
			return 1
		else
			echo "${WORKSPACE_NAME} started frontend_proc"
			return 0
		fi
	else
		exitcode=$?
		echo "${WORKSPACE_NAME} aborting frontend_proc due to exitcode ${exitcode} from frontend_proc"
		return $exitcode
	fi
}


run_all() {
	echo "Running frontend_ctr"

	# Check that all necessary environment variables are set
	echo "Required environment variables:"
	missing_vars=0
	if [ -z "${FRONTEND_SERVICE_HTTP_BIND_ADDR+x}" ]; then
		echo "  FRONTEND_SERVICE_HTTP_BIND_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  FRONTEND_SERVICE_HTTP_BIND_ADDR=$FRONTEND_SERVICE_HTTP_BIND_ADDR"
	fi
	
	if [ -z "${LUGGAGE_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "  LUGGAGE_SERVICE_HTTP_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  LUGGAGE_SERVICE_HTTP_DIAL_ADDR=$LUGGAGE_SERVICE_HTTP_DIAL_ADDR"
	fi
	
	if [ -z "${RESERV_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "  RESERV_SERVICE_HTTP_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  RESERV_SERVICE_HTTP_DIAL_ADDR=$RESERV_SERVICE_HTTP_DIAL_ADDR"
	fi
	
	if [ -z "${REVIEW_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "  REVIEW_SERVICE_HTTP_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  REVIEW_SERVICE_HTTP_DIAL_ADDR=$REVIEW_SERVICE_HTTP_DIAL_ADDR"
	fi
	
	if [ -z "${SEARCH_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "  SEARCH_SERVICE_HTTP_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  SEARCH_SERVICE_HTTP_DIAL_ADDR=$SEARCH_SERVICE_HTTP_DIAL_ADDR"
	fi
	
	if [ -z "${USER_SERVICE_HTTP_DIAL_ADDR+x}" ]; then
		echo "  USER_SERVICE_HTTP_DIAL_ADDR (missing)"
		missing_vars=$((missing_vars+1))
	else
		echo "  USER_SERVICE_HTTP_DIAL_ADDR=$USER_SERVICE_HTTP_DIAL_ADDR"
	fi
		

	if [ "$missing_vars" -gt 0 ]; then
		echo "Aborting due to missing environment variables"
		return 1
	fi

	frontend_proc
	
	wait
}

run_all
