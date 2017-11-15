package errors

import "github.com/normegil/resterrors"

const DEFAULT_ERROR_CODE = 50000

func Definitions() []resterrors.ErrorDefinition {
	baseURL := "http://example.com/rest/errors/"
	return []resterrors.ErrorDefinition{
		{50000, 500, baseURL + "50000", "An undetermined error happened on the Server."},
		{50301, 503, baseURL + "50301", "Could not connect to Zookeeper at given address. Check that zookeeper is running and accessible by the rest server."},
		{40001, 400, baseURL + "40001", "The request doesn't correspond to the structure needed to solve your request. Please review the body of your request."},
		{40002, 400, baseURL + "40002", "A value is missing from your request or is misplaced."},
		{40003, 400, baseURL + "40003", "One of your parameter doesn't have the expected format and cannot be parsed."},
		{40004, 400, baseURL + "40004", "Cannot remove root path from zookeeper."},
		{40005, 400, baseURL + "40005", "Cannot remove a node with existing childs (Use 'recursive' option)."},
		{40100, 401, baseURL + "40100", "All authentication headers were detected as empty."},
		{40101, 401, baseURL + "40101", "Error while trying to authenticate user using Authorization header."},
		{40102, 401, baseURL + "40102", "No error happened but user is empty."},
		{40103, 401, baseURL + "40103", "Authentication header didn't contain any user."},
	}
}
