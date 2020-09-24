package bob

import (
    "regexp"
)

/*
Bob answers 'Sure.' if you ask him a question, such as "How are you?".
He answers 'Whoa, chill out!' if you YELL AT HIM (in all capitals).
He answers 'Calm down, I know what I'm doing!' if you yell a question at him.
He says 'Fine. Be that way!' if you address him without actually saying anything.
He answers 'Whatever.' to anything else.
*/

const respSure string = "Sure."
const respWhoa string = "Whoa, chill out!"
const respCalm string = "Calm down, I know what I'm doing!"
const respFine string = "Fine. Be that way!"
const respWhtw string = "Whatever."


var patternQuestion = regexp.MustCompile(`\?+\W*$`)
var patternAllCaps = regexp.MustCompile("^(([0-9]*\\W*)*([A-Z]+\\W*)+([0-9]*\\W*)*)$")
var patternAddress = regexp.MustCompile(`^\W*$`)

func Hey(remark string) string {

    isQuestion := patternQuestion.MatchString(remark)
    isYelling := patternAllCaps.MatchString(remark)
    isAddress := patternAddress.MatchString(remark)

    switch {
        case isQuestion && isYelling:
            return respCalm
        case isQuestion:
            return respSure
        case isYelling:
            return respWhoa
        case isAddress:
            return respFine
        default:
            return respWhtw
    }
}
