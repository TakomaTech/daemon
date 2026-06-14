extern "C" int processSamples(float* samples, int len) {
    for (int i = 0; i < len; ++i) {
        samples[i] *= 1.0f;
    }
    return 0;
}
