package utils
import(
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func GenHash(pwd string, saltRounds int)(string,error){
	hash,err:=bcrypt.GenerateFromPassword([]byte(pwd),saltRounds)
	if err!= nil{
		return "",err
	}
	return string(hash),nil
}

func CheckPassword(pwd string, hash string)(bool,error){
	err:= bcrypt.CompareHashAndPassword([]byte(hash),[]byte(pwd))
	if err!=nil{
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("error comparing password with hash: %v",err)
	}
	return true, nil
}