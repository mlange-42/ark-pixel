# Ark Pixel

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/ark-pixel/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/ark-pixel/actions/workflows/tests.yml)
[![Coverage Status](https://badge.coveralls.io/repos/github/mlange-42/ark-pixel/badge.svg?branch=main)](https://badge.coveralls.io/github/mlange-42/ark-pixel?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/ark-pixel)](https://goreportcard.com/report/github.com/mlange-42/ark-pixel)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/ark-pixel.svg)](https://pkg.go.dev/github.com/mlange-42/ark-pixel)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/ark-pixel)
[![MIT license](https://img.shields.io/github/license/mlange-42/ark-pixel)](https://github.com/mlange-42/ark-pixel/blob/main/LICENSE)

*Ark Pixel* provides OpenGL graphics and live plots for the [Ark](https://github.com/mlange-42/ark) Entity Component System (ECS) using the [Pixel](https://github.com/gopxl/pixel) game engine.

<div align="center">

<a href="https://github.com/mlange-42/ark">
<img src="https://github.com/user-attachments/assets/4bbe57c6-2e16-43be-ad5e-0cf26c220f21" alt="Ark (logo)" width="500px" />
</a>

</div>
</br>
<div align="center" width="100%">

![Screenshot](https://user-images.githubusercontent.com/44003176/232126308-60299642-0490-478d-82a5-48d862da6703.png)  
*Screenshot showing Ark Pixel features, visualizing an evolutionary forest model.*
</div>

## Features

* Free 2D drawing using a convenient OpenGL interface.
* Live plots using unified observers (time series, line, bar, scatter and contour plots).
* ECS engine monitor for detailed performance statistics.
* Entity inspector for debugging and inspection.
* Simulation controls to pause or limit speed interactively.
* User input handling for interactive simulations.

## Installation

```
go get github.com/mlange-42/ark-pixel
```

The **dependencies** of [go-gl/gl](https://github.com/go-gl/gl) and [go-gl/glfw](https://github.com/go-gl/glfw) apply:

- A cgo compiler (typically gcc).
- For Ubuntu/Debian-based systems, you also need `libgl1-mesa-dev` and `xorg-dev`

## Usage

See the [API docs](https://pkg.go.dev/github.com/mlange-42/ark-pixel) for details and examples.

[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/ark-pixel.svg)](https://pkg.go.dev/github.com/mlange-42/ark-pixel)

## License

This project is distributed under the [MIT licence](./LICENSE).
