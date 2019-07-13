package test.projo;

import java.lang.annotation.Annotation;
import java.lang.annotation.IncompleteAnnotationException;

public class TestAnnotation extends IncompleteAnnotationException implements Annotation {

    public TestAnnotation(Class<? extends Annotation> annotationType, String elementName) {
        super(annotationType, elementName);
    }

    @Override
    public String elementName() {
        return IncompleteAnnotationException.class.getName();
    }

    @Override
    public Class<? extends Annotation> annotationType() {
        return TestAnnotation.class;
    }
}
