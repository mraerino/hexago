package main

import (
    "testing"
    "reflect"
)

func TestHexagonCalc(t *testing.T) {
    h_in := Hexagon{
        M: Point{X: 250, Y: 250,},
        R: 200,
    }
    h_out := Hexagon{
        A: Point{X: 250, Y: 50},
        B: Point{X: 423, Y: 150},
        C: Point{X: 423, Y: 350},
        D: Point{X: 250, Y: 450},
        E: Point{X: 77, Y: 350},
        F: Point{X: 77, Y: 150},
        M: h_in.M,
        R: h_in.R,
    }
    h_in.calc()
    if !reflect.DeepEqual(h_in, h_out) {
        t.Errorf("Hexagon not equal:\ngot\t\t%+v,\nexpected\t%+v", h_in, h_out)
    }
}
