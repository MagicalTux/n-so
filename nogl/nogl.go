package nogl

// OpenGL based UI for n-so

import (
	"log"
	"runtime"

	// OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Context struct {
	cl chan struct{}

	window  *glfw.Window
	program uint32
}

func New() (*Context, error) {
	ctx := &Context{
		cl: make(chan struct{}),
	}

	ch := make(chan error)
	go ctx.run(ch)

	return ctx, <-ch
}

func (c *Context) Wait() {
	<-c.cl
}

func (c *Context) run(ch chan error) {
	runtime.LockOSThread()

	err := c.initGlfw()
	if err != nil {
		ch <- err
		return
	}
	defer glfw.Terminate()

	err = c.initOpenGL()
	if err != nil {
		ch <- err
		return
	}

	// init done, let main process return
	ch <- nil

	for !c.window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(c.program)

		c.window.SwapBuffers()
		glfw.PollEvents()
	}
	close(c.cl)
}

// initGlfw initializes glfw and returns a Window to use.
func (c *Context) initGlfw() error {
	err := glfw.Init()
	if err != nil {
		return err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	//glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	c.window, err = glfw.CreateWindow(640, 480, "n-so", nil, nil)
	if err != nil {
		glfw.Terminate() // won't be called by parent
		return err
	}
	c.window.MakeContextCurrent()

	return nil
}

func (c *Context) initOpenGL() error {
	if err := gl.Init(); err != nil {
		return err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	c.program = gl.CreateProgram()
	gl.LinkProgram(c.program)
	return nil
}
