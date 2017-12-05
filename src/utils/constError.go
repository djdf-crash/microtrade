package utils

var UserRegisterError = map[int]string{
	101: "User with this email is already exists",
	102: "The password did not pass the complexity check or invalid",
	103: "Device id contains invalid characters",
}

var PasswordResetRequestError = map[int]string{
	301: "User with such email not found",
}

var PasswordResetConfirmationError = map[int]string{
	401: "User uid for reset password is invalid",
	402: "Token for reset password is invalid",
}

var ChangeUserPasswordError = map[int]string{
	601: "Current password is invalid",
	602: "New password is invalid",
}

var ValidationReqError = map[int]string{
	1:  "Field is required",
	2:  "Field cannot be null",
	3:  "Field cannot be blank",
	4:  "Field is invalid",
	5:  "Field contains too many characters",
	6:  "Field contains less characters than it is needed",
	7:  "Field contains greater value than it is needed",
	8:  "Field contains less value than it is needed",
	9:  "Field contains value too large (number decoded to string)",
	10: "The total number of digits in the field is greater than allowed",
	11: "The decimal places of digits in the field is greater than allowed",
	12: "The digits before the decimal point in the field is greater than allowed",
	13: "Expected a datetime but got a date",
	14: "Invalid datetime for the used timezone",
	15: "Expected a date but got a datetime",
	16: "Field contains is not a valid choice",
	17: "Field contains is not a list of items",
	18: "Field may not be empty",
}

var CommonError = map[int]string{
	-1: "Access Denied",
	-2: "Error parsing input JSON",
	-3: "Token invalid or expired",
}
