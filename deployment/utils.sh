#!/bin/bash


new_tag ()
{
   date +%Y%m%d"."${CIRCLE_BUILD_NUM}
}

