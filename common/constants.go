// Copyright © 2018 SUSE LLC
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

package common

const STATE_RUNNING = "running"

const AssignSep = "@"
const AssignFmt = "%s" + AssignSep + "%s" + AssignSep + "%s"

const StateSep = "!"
const StateFmt = "%s" + StateSep + "%s"

const WorkerSep = "#"
const WorkerInstSep = ":"
const WorkerClassSep = ","
const WorkerEncodeFormat = "%s" + WorkerInstSep + "%d" + WorkerSep + "%s"

const TestSep = "#"
const TestParallelSep = ","
const TestEncodeFormat = "%s" + TestSep + "%s" + TestSep + "%s" + TestSep + "%s"

// sched states
const STATE_OLD = "old"
const STATE_CURRENT = "current"
