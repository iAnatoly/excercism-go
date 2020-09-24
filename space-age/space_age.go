package space

type Planet string

var xlat = map[Planet]float64 {
    "Mercury": 0.2408467,
    "Venus": 0.61519726,
    "Earth": 1.0,
    "Mars": 1.8808158,
    "Jupiter": 11.862615,
    "Saturn": 29.447498,
    "Uranus": 84.016846,
    "Neptune": 164.79132,
}

const seconds_in_year float64 = 31557600.0

func Age(age_in_seconds float64, planet Planet) float64 {
    var age_in_earth_years = age_in_seconds / seconds_in_year
    return age_in_earth_years / xlat[planet]
}
