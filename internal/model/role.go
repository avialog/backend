package model

type Role string

const (
	RolePilotInCommand                 Role = "PIC"
	RoleSecondInCommand                Role = "SIC"
	RoleDual                           Role = "DUAL"
	RoleStudentPilotInCommand          Role = "SPIC"
	RolePilotInCommandUnderSupervision Role = "P1S"
	RoleInstructor                     Role = "INS"
	RoleExaminer                       Role = "EXM"
	RoleFlightAttendant                Role = "ATT"
	RoleOther                          Role = "OTH"
)
