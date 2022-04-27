# gof1

This library is to be use to communicate and control with a [Native Instrument
F1](https://www.native-instruments.com/en/products/traktor/dj-controllers/traktor-kontrol-f1/).

F1 controller aren't MIDI device per se, they are HID devices just like your
mouse or keyboard. The NI drivers converts signal to/from the device into MIDI
signals (which are set through the Controller Editor app).

This library implemented the HID protocol and provide some utilities to setup
control layouts and states.
