all:
	go run main.go draw.go matrix.go line.go transform.go parse.go parametric.go shapes.go stack.go screen.go lighting.go
	#animate rolling.gif
	#animate halo2.gif
clean:
	touch fo1232121o.png
	touch fo23231o.ppm
	touch parse/f12312oo23.pyc
	touch 123122.out
	touch img/321431241.png
	touch qjirjqijri2.gif
	rm *.gif
	rm *.png
	rm *.ppm
	rm parse/*.pyc
	rm *.out
	rm img/*.png
